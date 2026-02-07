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

func (r *Runner) Build(ctx context.Context) (*FakeTime, error) {
	s, err := r.provider.Get(ctx)
	if err != nil {
		return nil, err
	}

	fakeTime, err := Parse(s, r.layout)
	if err != nil {
		return nil, err
	}

	return fakeTime, nil
}

func (r *Runner) Start(ctx context.Context, fn func(context.Context) error) error {
	fakeTime, err := r.Build(ctx)
	if err != nil {
		return err
	}
	defer fakeTime.Setup(ctx)()
	return fn(ctx)
}
