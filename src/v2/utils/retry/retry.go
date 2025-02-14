package retry

import (
	"context"
	"log/slog"
	"math"
	"math/rand"
	"time"
)

type Func func() error
type FuncWithData[T any] func() (T, error)
type Option func(*Options)

type Options struct {
	MaxAttempts uint
	IsDelay     bool //是否开启延迟
	Delay       *Delay
}

func DefaultOptions() *Options {
	return &Options{
		MaxAttempts: 3,
		IsDelay:     false,
		Delay:       defaultDelay(),
	}
}

func defaultDelay() *Delay {
	return &Delay{
		Min:    500 * time.Millisecond,
		Max:    2 * time.Minute,
		Jitter: 0.2,
	}
}

type Delay struct {
	Min    time.Duration
	Max    time.Duration
	Jitter float64
}

func WithMaxAttempts(maxAttempts uint) Option {
	return func(o *Options) {
		o.MaxAttempts = maxAttempts
	}
}

func WithIsDelay(isDelay bool) Option {
	return func(o *Options) {
		o.IsDelay = isDelay
	}
}
func WithDelay(min, max time.Duration, jitter float64) Option {
	return func(o *Options) {
		o.Delay = &Delay{
			Min:    min,
			Max:    max,
			Jitter: jitter,
		}
	}
}

func DelayTime(delay, maxDelay time.Duration, jitter float64, attempt uint) time.Duration {
	delay = time.Duration(math.Pow(2, float64(attempt))) * delay
	if delay > maxDelay {
		delay = maxDelay
	}
	if jitter > 0 {
		delay += time.Duration(rand.Float64() * (jitter * float64(delay)))
	}
	return delay

}

func Retry(ctx context.Context, f Func, opts ...Option) error {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	var err error
	var attempt uint
	for {
		// MaxAttempts==0 即进入死循环,直到成功为止
		if options.MaxAttempts > 0 && attempt >= options.MaxAttempts {
			break
		}
		if err = f(); err == nil {
			break
		}
		var duration time.Duration
		if options.IsDelay {
			duration = DelayTime(options.Delay.Min, options.Delay.Max, options.Delay.Jitter, attempt)
		}
		select {
		case <-time.After(duration):
		case <-ctx.Done():
			return ctx.Err()
		}
		attempt++
		slog.Info("Retry", "options", options, "attempt", attempt, "duration", duration.Seconds(), "err", err)
	}
	return err
}

// RetryWithData 重试
func RetryWithData[T any](ctx context.Context, f FuncWithData[T], opts ...Option) (T, error) {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	var err error
	var attempt uint
	var data T
	for {
		data, err = f()
		if err == nil {
			break
		}
		// MaxAttempts==0 即进入死循环,直到成功为止
		if options.MaxAttempts > 0 && attempt >= options.MaxAttempts {
			break
		}
		var duration time.Duration
		if options.IsDelay {
			duration = DelayTime(options.Delay.Min, options.Delay.Max, options.Delay.Jitter, attempt)
		}
		select {
		case <-time.After(duration):
		case <-ctx.Done():
			return data, ctx.Err()
		}
		attempt++
		slog.Info("RetryWithData", "options", options, "attempt", attempt, "duration", duration.Seconds(), "err", err)
	}
	return data, err
}
