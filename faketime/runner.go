package faketime

import (
	"context"
	"errors"
	orig "time"

	"github.com/akm/time"
	"github.com/akm/time/testtime"
)

type Runner struct {
	provider Provider
	layout   string
}

func NewRunner(provider Provider, layout string) *Runner {
	return &Runner{
		provider: provider,
		layout:   layout,
	}
}

var (
	ErrInvalidFaketimeFileContent = errors.New("invalid faketime file content")
)

func (r *Runner) Start(ctx context.Context, fn func(context.Context) error) error {
	s, err := r.provider.Get(ctx)
	if err != nil {
		return err
	}

	fakeTime, err := Parse(s, r.layout)
	if err != nil {
		return err
	}

	if fakeTime.Ratio == 0 {
		defer testtime.SetTime(&fakeTime.Time)()
		return fn(ctx)
	}

	t0 := orig.Now()
	defer testtime.SetTimeFunc(func() time.Time {
		elapsed := time.Duration(float64(orig.Since(t0)) * fakeTime.Ratio)
		return fakeTime.Time.Add(elapsed)
	})()
	return fn(ctx)
}
