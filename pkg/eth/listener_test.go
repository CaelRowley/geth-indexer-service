package eth

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

type fakeSubscription struct {
	errCh        chan error
	unsubscribed atomic.Bool
}

func (f *fakeSubscription) Unsubscribe() {
	f.unsubscribed.Store(true)
}

func (f *fakeSubscription) Err() <-chan error {
	return f.errCh
}

type failingSubscriber struct {
	err error
}

func (f *failingSubscriber) SubscribeNewHead(_ context.Context, _ chan<- *types.Header) (ethereum.Subscription, error) {
	return nil, f.err
}

func TestRunListenerBackoff(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	subscriber := &failingSubscriber{err: errors.New("boom")}
	var sleepCalls []time.Duration

	err := runListener(ctx, subscriber, func(context.Context, *types.Header) error {
		return nil
	}, listenOptions{
		initialBackoff: time.Second,
		maxBackoff:     time.Second * 5,
		sleep: func(ctx context.Context, duration time.Duration) error {
			sleepCalls = append(sleepCalls, duration)
			if len(sleepCalls) == 2 {
				cancel()
				return ctx.Err()
			}
			return nil
		},
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(sleepCalls) != 2 {
		t.Fatalf("expected 2 sleep calls, got %d", len(sleepCalls))
	}
	if sleepCalls[0] != time.Second || sleepCalls[1] != time.Second*2 {
		t.Fatalf("unexpected backoff sequence: %v", sleepCalls)
	}
}

func TestConsumeHeadersStopsOnContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sub := &fakeSubscription{errCh: make(chan error)}
	headerCh := make(chan *types.Header, 1)
	var handled atomic.Int32

	go func() {
		headerCh <- &types.Header{}
	}()

	handler := func(_ context.Context, _ *types.Header) error {
		handled.Add(1)
		cancel()
		return nil
	}

	if err := consumeHeaders(ctx, sub, headerCh, handler); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if handled.Load() != 1 {
		t.Fatalf("expected handler to be called once, got %d", handled.Load())
	}
	if !sub.unsubscribed.Load() {
		t.Fatalf("expected subscription to be unsubscribed")
	}
}
