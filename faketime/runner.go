package faketime

import (
	"context"
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

func (r *Runner) Start(ctx context.Context, fn func(context.Context) error) error {
	s, err := r.provider.Get(ctx)
	if err != nil {
		return err
	}

	fakeTime, err := Parse(s, r.layout)
	if err != nil {
		return err
	}

	return fakeTime.Run(ctx, fn)
}
