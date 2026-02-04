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

func TestRunListenerStopsOnContextBeforeSubscribe(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	subscriber := &failingSubscriber{err: errors.New("should not be called")}
	var sleepCalled bool

	err := runListener(ctx, subscriber, func(context.Context, *types.Header) error {
		t.Fatal("handler should not be called")
		return nil
	}, listenOptions{
		initialBackoff: time.Second,
		maxBackoff:     time.Second * 5,
		sleep: func(ctx context.Context, duration time.Duration) error {
			sleepCalled = true
			return ctx.Err()
		},
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if sleepCalled {
		t.Fatalf("sleep should not be called when context is cancelled")
	}
}

func TestRunListenerMaxBackoff(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	subscriber := &failingSubscriber{err: errors.New("boom")}
	var sleepCalls []time.Duration

	err := runListener(ctx, subscriber, func(context.Context, *types.Header) error {
		return nil
	}, listenOptions{
		initialBackoff: time.Second,
		maxBackoff:     time.Second * 4,
		sleep: func(ctx context.Context, duration time.Duration) error {
			sleepCalls = append(sleepCalls, duration)
			if len(sleepCalls) == 4 {
				cancel()
				return ctx.Err()
			}
			return nil
		},
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	expected := []time.Duration{time.Second, time.Second * 2, time.Second * 4, time.Second * 4}
	if len(sleepCalls) != len(expected) {
		t.Fatalf("expected %d sleep calls, got %d", len(expected), len(sleepCalls))
	}
	for i, exp := range expected {
		if sleepCalls[i] != exp {
			t.Fatalf("expected sleepCalls[%d] to be %v, got %v", i, exp, sleepCalls[i])
		}
	}
}

func TestRunListenerResetsBackoffAfterSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	attemptCount := 0
	subscriber := &retryingSubscriber{
		failCount: 2,
		attemptCount: &attemptCount,
		errCh: make(chan error),
	}
	var sleepCalls []time.Duration

	go func() {
		time.Sleep(10 * time.Millisecond)
		subscriber.errCh <- errors.New("connection dropped")
	}()

	err := runListener(ctx, subscriber, func(context.Context, *types.Header) error {
		return nil
	}, listenOptions{
		initialBackoff: time.Millisecond * 10,
		maxBackoff:     time.Second,
		sleep: func(ctx context.Context, duration time.Duration) error {
			sleepCalls = append(sleepCalls, duration)
			if len(sleepCalls) == 3 {
				cancel()
				return ctx.Err()
			}
			return nil
		},
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(sleepCalls) < 3 {
		t.Fatalf("expected at least 3 sleep calls, got %d", len(sleepCalls))
	}
	if sleepCalls[0] != time.Millisecond*10 || sleepCalls[1] != time.Millisecond*20 {
		t.Fatalf("unexpected initial backoff sequence: %v", sleepCalls[:2])
	}
	if sleepCalls[2] != time.Millisecond*10 {
		t.Fatalf("expected backoff to reset to initial value after successful subscription, got %v", sleepCalls[2])
	}
}

func TestRunListenerDefaultOptions(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	subscriber := &failingSubscriber{err: errors.New("boom")}
	var sleepCalled bool

	err := runListener(ctx, subscriber, func(context.Context, *types.Header) error {
		return nil
	}, listenOptions{
		sleep: func(ctx context.Context, duration time.Duration) error {
			sleepCalled = true
			if duration != time.Second {
				t.Fatalf("expected default initialBackoff of 1s, got %v", duration)
			}
			cancel()
			return ctx.Err()
		},
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !sleepCalled {
		t.Fatalf("expected sleep to be called")
	}
}

func TestRunListenerNilSleepUsesDefault(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	subscriber := &failingSubscriber{err: errors.New("boom")}
	attemptCount := 0

	go func() {
		time.Sleep(5 * time.Millisecond)
		cancel()
	}()

	err := runListener(ctx, subscriber, func(context.Context, *types.Header) error {
		attemptCount++
		return nil
	}, listenOptions{
		initialBackoff: time.Millisecond,
		maxBackoff:     time.Millisecond * 5,
		sleep:          nil,
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if attemptCount > 0 {
		t.Fatalf("handler should not have been called when subscription fails")
	}
}

func TestConsumeHeadersSubscriptionError(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("subscription error")
	sub := &fakeSubscription{errCh: make(chan error, 1)}
	headerCh := make(chan *types.Header)

	sub.errCh <- expectedErr

	err := consumeHeaders(ctx, sub, headerCh, func(context.Context, *types.Header) error {
		t.Fatal("handler should not be called")
		return nil
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
	if !sub.unsubscribed.Load() {
		t.Fatalf("expected subscription to be unsubscribed")
	}
}

func TestConsumeHeadersSubscriptionClosed(t *testing.T) {
	ctx := context.Background()

	sub := &fakeSubscription{errCh: make(chan error)}
	headerCh := make(chan *types.Header)

	close(sub.errCh)

	err := consumeHeaders(ctx, sub, headerCh, func(context.Context, *types.Header) error {
		t.Fatal("handler should not be called")
		return nil
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "subscription closed" {
		t.Fatalf("expected 'subscription closed' error, got %v", err)
	}
	if !sub.unsubscribed.Load() {
		t.Fatalf("expected subscription to be unsubscribed")
	}
}

func TestConsumeHeadersNilHeaderIgnored(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	sub := &fakeSubscription{errCh: make(chan error)}
	headerCh := make(chan *types.Header, 2)
	var handled atomic.Int32

	headerCh <- nil
	headerCh <- &types.Header{}

	handler := func(_ context.Context, _ *types.Header) error {
		handled.Add(1)
		cancel()
		return nil
	}

	err := consumeHeaders(ctx, sub, headerCh, handler)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if handled.Load() != 1 {
		t.Fatalf("expected handler to be called once, got %d", handled.Load())
	}
}

func TestConsumeHeadersHandlerError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sub := &fakeSubscription{errCh: make(chan error)}
	headerCh := make(chan *types.Header, 2)
	var callCount atomic.Int32

	headerCh <- &types.Header{}
	headerCh <- &types.Header{}

	handler := func(_ context.Context, _ *types.Header) error {
		count := callCount.Add(1)
		if count == 1 {
			return errors.New("handler error")
		}
		cancel()
		return nil
	}

	err := consumeHeaders(ctx, sub, headerCh, handler)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if callCount.Load() != 2 {
		t.Fatalf("expected handler to be called twice, got %d", callCount.Load())
	}
}

func TestSleepWithContextZeroDuration(t *testing.T) {
	ctx := context.Background()
	err := sleepWithContext(ctx, 0)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestSleepWithContextNegativeDuration(t *testing.T) {
	ctx := context.Background()
	err := sleepWithContext(ctx, -time.Second)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestSleepWithContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := sleepWithContext(ctx, time.Second)
	if err != context.Canceled {
		t.Fatalf("expected context.Canceled error, got %v", err)
	}
}

func TestSleepWithContextCancelledDuringSleep(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() {
		done <- sleepWithContext(ctx, time.Second)
	}()

	time.Sleep(10 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err != context.Canceled {
			t.Fatalf("expected context.Canceled error, got %v", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("sleepWithContext did not return after context cancellation")
	}
}

func TestSleepWithContextCompletes(t *testing.T) {
	ctx := context.Background()
	start := time.Now()
	err := sleepWithContext(ctx, 50*time.Millisecond)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if elapsed < 50*time.Millisecond {
		t.Fatalf("sleep completed too early: %v", elapsed)
	}
}

func TestNextBackoffDoubles(t *testing.T) {
	tests := []struct {
		current  time.Duration
		max      time.Duration
		expected time.Duration
	}{
		{time.Second, time.Second * 10, time.Second * 2},
		{time.Second * 2, time.Second * 10, time.Second * 4},
		{time.Millisecond * 500, time.Second * 10, time.Second},
		{time.Millisecond * 100, time.Second, time.Millisecond * 200},
	}

	for _, tt := range tests {
		result := nextBackoff(tt.current, tt.max)
		if result != tt.expected {
			t.Errorf("nextBackoff(%v, %v) = %v, want %v", tt.current, tt.max, result, tt.expected)
		}
	}
}

func TestNextBackoffCapsAtMax(t *testing.T) {
	tests := []struct {
		current  time.Duration
		max      time.Duration
		expected time.Duration
	}{
		{time.Second * 5, time.Second * 5, time.Second * 5},
		{time.Second * 6, time.Second * 10, time.Second * 10},
		{time.Second * 30, time.Second * 30, time.Second * 30},
		{time.Second * 100, time.Second * 50, time.Second * 50},
	}

	for _, tt := range tests {
		result := nextBackoff(tt.current, tt.max)
		if result != tt.expected {
			t.Errorf("nextBackoff(%v, %v) = %v, want %v", tt.current, tt.max, result, tt.expected)
		}
	}
}

type retryingSubscriber struct {
	failCount    int
	attemptCount *int
	errCh        chan error
}

func (r *retryingSubscriber) SubscribeNewHead(_ context.Context, _ chan<- *types.Header) (ethereum.Subscription, error) {
	*r.attemptCount++
	if *r.attemptCount <= r.failCount {
		return nil, errors.New("subscription failed")
	}
	return &fakeSubscription{errCh: r.errCh}, nil
}

type successfulSubscriber struct {
	headerCh chan *types.Header
	sub      *fakeSubscription
}

func (s *successfulSubscriber) SubscribeNewHead(_ context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	go func() {
		for header := range s.headerCh {
			ch <- header
		}
	}()
	return s.sub, nil
}

func TestRunListenerSuccessfulFlow(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	headerCh := make(chan *types.Header, 1)
	sub := &fakeSubscription{errCh: make(chan error)}
	subscriber := &successfulSubscriber{
		headerCh: headerCh,
		sub:      sub,
	}

	var handledHeaders atomic.Int32

	go func() {
		time.Sleep(10 * time.Millisecond)
		headerCh <- &types.Header{}
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	err := runListener(ctx, subscriber, func(_ context.Context, _ *types.Header) error {
		handledHeaders.Add(1)
		return nil
	}, listenOptions{
		initialBackoff: time.Millisecond,
		maxBackoff:     time.Second,
		sleep:          sleepWithContext,
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if handledHeaders.Load() != 1 {
		t.Fatalf("expected 1 header to be handled, got %d", handledHeaders.Load())
	}
	if !sub.unsubscribed.Load() {
		t.Fatal("expected subscription to be unsubscribed")
	}
}