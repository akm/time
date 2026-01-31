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

type FakeTime struct {
	Time  time.Time
	Ratio float64
}

var (
	ErrInvalidFaketimeFileContent = errors.New("invalid faketime file content")
)

func Parse(s string, layout string) (*FakeTime, error) {
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
			return nil, fmt.Errorf("%w: failed to parse duration from file content: %s, error: %v", ErrInvalidFaketimeFileContent, s, err)
		}
		t = time.Now().Add(offsetDuration)
	} else {
		body = strings.TrimPrefix(body, "@")
		var err error
		t, err = time.Parse(layout, body)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to parse time from file content: %s, error: %v", ErrInvalidFaketimeFileContent, s, err)
		}
	}

	if opt == "" {
		return &FakeTime{Time: t, Ratio: 0}, nil
	}

	var ratioStr string
	if strings.HasPrefix(opt, "x") {
		ratioStr = strings.TrimPrefix(opt, "x")
		if ratioStr == "" {
			return nil, fmt.Errorf("%w: missing ratio after 'x' in file content: %s", ErrInvalidFaketimeFileContent, s)
		}
	}

	var ratio float64
	if ratioStr == "" {
		ratio = 1.0
	} else {
		var err error
		ratio, err = strconv.ParseFloat(ratioStr, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to parse ratio from file content: %s, error: %v", ErrInvalidFaketimeFileContent, s, err)
		}
		if ratio <= 0 {
			return nil, fmt.Errorf("%w: ratio must be positive in file content: %s", ErrInvalidFaketimeFileContent, s)
		}
	}

	return &FakeTime{Time: t, Ratio: ratio}, nil
}

func (ft *FakeTime) Run(ctx context.Context, fn func(context.Context) error) error {
	if ft.Ratio == 0 {
		defer testtime.SetTime(&ft.Time)()
		return fn(ctx)
	}
	t0 := orig.Now()
	defer testtime.SetTimeFunc(func() time.Time {
		elapsed := time.Duration(float64(orig.Since(t0)) * ft.Ratio)
		return ft.Time.Add(elapsed)
	})()
	return fn(ctx)
}
