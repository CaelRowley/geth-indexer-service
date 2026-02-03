package eth

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/CaelRowley/geth-indexer-service/pkg/db"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

type headSubscriber interface {
	SubscribeNewHead(context.Context, chan<- *types.Header) (ethereum.Subscription, error)
}

type headerHandler func(context.Context, *types.Header) error

type listenOptions struct {
	initialBackoff time.Duration
	maxBackoff     time.Duration
	sleep          func(context.Context, time.Duration) error
}

func (c EthClient) StartListener(ctx context.Context, dbConn db.DB) error {
	return runListener(ctx, c, c.handleHeader, listenOptions{
		initialBackoff: time.Second,
		maxBackoff:     time.Second * 30,
		sleep:          sleepWithContext,
	})
}

func runListener(ctx context.Context, subscriber headSubscriber, handler headerHandler, opts listenOptions) error {
	if opts.sleep == nil {
		opts.sleep = sleepWithContext
	}
	if opts.initialBackoff <= 0 {
		opts.initialBackoff = time.Second
	}
	if opts.maxBackoff <= 0 {
		opts.maxBackoff = time.Second * 30
	}

	headerCh := make(chan *types.Header)
	defer close(headerCh)

	backoff := opts.initialBackoff
	for {
		if ctx.Err() != nil {
			slog.Info("eth listener stopped")
			return nil
		}

		sub, err := subscriber.SubscribeNewHead(ctx, headerCh)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			slog.Warn("failed to subscribe to head; retrying", "err", err, "backoff", backoff)
			if err := opts.sleep(ctx, backoff); err != nil {
				return nil
			}
			backoff = nextBackoff(backoff, opts.maxBackoff)
			continue
		}
		backoff = opts.initialBackoff

		if err := consumeHeaders(ctx, sub, headerCh, handler); err != nil {
			if ctx.Err() != nil {
				return nil
			}
			slog.Warn("subscription dropped; retrying", "err", err, "backoff", backoff)
			if err := opts.sleep(ctx, backoff); err != nil {
				return nil
			}
			backoff = nextBackoff(backoff, opts.maxBackoff)
			continue
		}
		return nil
	}
}

func consumeHeaders(ctx context.Context, sub ethereum.Subscription, headerCh <-chan *types.Header, handler headerHandler) error {
	defer sub.Unsubscribe()

	for {
		select {
		case err, ok := <-sub.Err():
			if !ok || err == nil {
				return fmt.Errorf("subscription closed")
			}
			return err
		case header := <-headerCh:
			if header == nil {
				continue
			}
			if err := handler(ctx, header); err != nil {
				slog.Error("failed to process header", "err", err)
			}
		case <-ctx.Done():
			slog.Info("eth listener stopped")
			return nil
		}
	}
}

func sleepWithContext(ctx context.Context, duration time.Duration) error {
	if duration <= 0 {
		return nil
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func nextBackoff(current, max time.Duration) time.Duration {
	next := current * 2
	if next > max {
		return max
	}
	return next
}

func (c EthClient) handleHeader(ctx context.Context, header *types.Header) error {
	block, err := c.BlockByHash(ctx, header.Hash())
	if err != nil {
		return err
	}
	if err = c.publishBlock(block); err != nil {
		return err
	}
	if err := c.publishTxs(ctx, block.Transactions(), block.Hash()); err != nil {
		return err
	}
	return nil
}
