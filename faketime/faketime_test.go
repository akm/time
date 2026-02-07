package faketime

import (
	"context"
	"errors"
	"testing"

	"github.com/akm/time"
	"github.com/akm/time/testtime"
)

func TestParse(t *testing.T) {
	layout := "2006-01-02 15:04:05"

	// Set a fixed "now" time for tests involving relative times
	baseTime := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	defer testtime.SetTime(&baseTime)()

	tests := []struct {
		name      string
		input     string
		wantTime  time.Time
		wantRatio float64
		wantErr   error
	}{
		{
			name:      "absolute time without prefix",
			input:     "2024-01-02 15:04:05",
			wantTime:  time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC),
			wantRatio: 0,
		},
		{
			name:      "absolute time with @ prefix",
			input:     "@2024-01-02 15:04:05",
			wantTime:  time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC),
			wantRatio: 0,
		},
		{
			name:      "absolute time with + suffix (ratio 1.0)",
			input:     "2024-01-02 15:04:05 +",
			wantTime:  time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC),
			wantRatio: 1.0,
		},
		{
			name:      "absolute time with x2 suffix (ratio 2.0)",
			input:     "2024-01-02 15:04:05 x2",
			wantTime:  time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC),
			wantRatio: 2.0,
		},
		{
			name:      "absolute time with x0.5 suffix (ratio 0.5)",
			input:     "2024-01-02 15:04:05 x0.5",
			wantTime:  time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC),
			wantRatio: 0.5,
		},
		{
			name:      "positive duration offset",
			input:     "+1h30m",
			wantTime:  baseTime.Add(time.Hour + 30*time.Minute),
			wantRatio: 0,
		},
		{
			name:      "negative duration offset",
			input:     "-30m",
			wantTime:  baseTime.Add(-30 * time.Minute),
			wantRatio: 0,
		},
		{
			name:      "positive duration with + suffix",
			input:     "+1h +",
			wantTime:  baseTime.Add(time.Hour),
			wantRatio: 1.0,
		},
		{
			name:      "positive duration with x suffix",
			input:     "+2h x3",
			wantTime:  baseTime.Add(2 * time.Hour),
			wantRatio: 3.0,
		},
		{
			name:    "invalid duration",
			input:   "+invalid",
			wantErr: ErrInvalidFaketimeFileContent,
		},
		{
			name:    "invalid time format",
			input:   "not-a-time",
			wantErr: ErrInvalidFaketimeFileContent,
		},
		{
			name:    "x without ratio",
			input:   "2024-01-02 15:04:05 x",
			wantErr: ErrInvalidFaketimeFileContent,
		},
		{
			name:    "x with invalid ratio",
			input:   "2024-01-02 15:04:05 xabc",
			wantErr: ErrInvalidFaketimeFileContent,
		},
		{
			name:    "x with zero ratio",
			input:   "2024-01-02 15:04:05 x0",
			wantErr: ErrInvalidFaketimeFileContent,
		},
		{
			name:    "x with negative ratio",
			input:   "2024-01-02 15:04:05 x-1",
			wantErr: ErrInvalidFaketimeFileContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input, layout)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !got.Time.Equal(tt.wantTime) {
				t.Errorf("Time = %v, want %v", got.Time, tt.wantTime)
			}
			if got.Ratio != tt.wantRatio {
				t.Errorf("Ratio = %v, want %v", got.Ratio, tt.wantRatio)
			}
		})
	}
}

func TestFakeTime_Setup(t *testing.T) {
	t.Run("ratio 0 sets fixed time", func(t *testing.T) {
		ft := FakeTime{
			Time:  time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC),
			Ratio: 0,
		}
		ctx := context.Background()

		cleanup := ft.Setup(ctx)
		defer cleanup()

		now := time.Now()
		expected := time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)
		if !now.Equal(expected) {
			t.Errorf("time.Now() = %v, want %v", now, expected)
		}
	})

	t.Run("ratio 1.0 starts time from specified point", func(t *testing.T) {
		ft := FakeTime{
			Time:  time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC),
			Ratio: 1.0,
		}
		ctx := context.Background()

		cleanup := ft.Setup(ctx)
		defer cleanup()

		start := time.Now()
		expected := time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)
		if start.Before(expected) {
			t.Errorf("time.Now() = %v, should be at or after %v", start, expected)
		}
	})

	t.Run("ratio 2.0 makes time pass twice as fast", func(t *testing.T) {
		ft := FakeTime{
			Time:  time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC),
			Ratio: 2.0,
		}
		ctx := context.Background()

		cleanup := ft.Setup(ctx)
		defer cleanup()

		start := time.Now()
		expected := time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)
		if start.Before(expected) {
			t.Errorf("time.Now() = %v, should be at or after %v", start, expected)
		}
	})

	t.Run("cleanup restores original time", func(t *testing.T) {
		ft := FakeTime{
			Time:  time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC),
			Ratio: 0,
		}
		ctx := context.Background()

		cleanup := ft.Setup(ctx)
		cleanup()

		now := time.Now()
		fakeTime := time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)
		if now.Equal(fakeTime) {
			t.Error("time should be restored after cleanup")
		}
	})
}

func TestFakeTime_Run(t *testing.T) {
	tests := []struct {
		name          string
		fakeTime      FakeTime
		checkTimeFunc func(t *testing.T)
	}{
		{
			name: "ratio 0 sets fixed time",
			fakeTime: FakeTime{
				Time:  time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC),
				Ratio: 0,
			},
			checkTimeFunc: func(t *testing.T) {
				now := time.Now()
				expected := time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)
				if !now.Equal(expected) {
					t.Errorf("time.Now() = %v, want %v", now, expected)
				}
			},
		},
		{
			name: "ratio 1.0 starts time from specified point",
			fakeTime: FakeTime{
				Time:  time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC),
				Ratio: 1.0,
			},
			checkTimeFunc: func(t *testing.T) {
				start := time.Now()
				expected := time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)
				// Time should be at or after the expected start time
				if start.Before(expected) {
					t.Errorf("time.Now() = %v, should be at or after %v", start, expected)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := tt.fakeTime.Run(ctx, func(ctx context.Context) error {
				tt.checkTimeFunc(t)
				return nil
			})

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestFakeTime_Run_ReturnsError(t *testing.T) {
	expectedErr := errors.New("test error")
	ft := FakeTime{
		Time:  time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC),
		Ratio: 0,
	}

	ctx := context.Background()
	err := ft.Run(ctx, func(ctx context.Context) error {
		return expectedErr
	})

	if err != expectedErr {
		t.Errorf("Run() error = %v, want %v", err, expectedErr)
	}
}
