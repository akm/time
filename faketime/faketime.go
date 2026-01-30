package faketime

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
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

func (ft *Runner) Start(ctx context.Context, fn func(context.Context) error) error {
	s, err := ft.provider.Get(ctx)
	if err != nil {
		return err
	}

	parts := strings.Split(s, " ")
	var body string
	var opt string
	if len(parts) > 1 {
		last := parts[len(parts)-1]
		if last == "+" || strings.HasPrefix(last, "x") {
			opt = last
			body = strings.Join(parts[:len(parts)-1], " ")
		} else {
			body = s
		}
	} else {
		body = s
	}

	var t time.Time
	if strings.HasPrefix(body, "+") || strings.HasPrefix(body, "-") {
		offsetDuration, err := time.ParseDuration(body)
		if err != nil {
			return fmt.Errorf("%w: failed to parse duration from file content: %s, error: %v", ErrInvalidFaketimeFileContent, s, err)
		}
		t = time.Now().Add(offsetDuration)
	} else {
		body = strings.TrimPrefix(body, "@")
		var err error
		t, err = time.Parse(ft.layout, body)
		if err != nil {
			return fmt.Errorf("%w: failed to parse time from file content: %s, error: %v", ErrInvalidFaketimeFileContent, s, err)
		}
	}

	if opt == "" {
		defer testtime.SetTime(&t)()
		return fn(ctx)
	}

	var ratioStr string
	if strings.HasPrefix(opt, "x") {
		ratioStr = strings.TrimPrefix(opt, "x")
		if ratioStr == "" {
			return fmt.Errorf("%w: missing ratio after 'x' in file content: %s", ErrInvalidFaketimeFileContent, s)
		}
	}

	var ratio float64
	if ratioStr == "" {
		ratio = 1.0
	} else {
		var err error
		ratio, err = strconv.ParseFloat(ratioStr, 64)
		if err != nil {
			return fmt.Errorf("%w: failed to parse ratio from file content: %s, error: %v", ErrInvalidFaketimeFileContent, s, err)
		}
		if ratio <= 0 {
			return fmt.Errorf("%w: ratio must be positive in file content: %s", ErrInvalidFaketimeFileContent, s)
		}
	}

	t0 := orig.Now()
	defer testtime.SetTimeFunc(func() time.Time {
		elapsed := time.Duration(float64(orig.Since(t0)) * ratio)
		return t.Add(elapsed)
	})()
	return fn(ctx)
}
