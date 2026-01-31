package faketime

import (
	"context"
	"errors"
	"testing"

	"github.com/akm/time"
)

// mockProvider implements Provider interface for testing
type mockProvider struct {
	value string
	err   error
}

func (m *mockProvider) Get(ctx context.Context) (string, error) {
	return m.value, m.err
}

var _ Provider = (*mockProvider)(nil)

func TestNewRunner(t *testing.T) {
	provider := &mockProvider{}
	layout := "2006-01-02 15:04:05"

	runner := NewRunner(provider, layout)

	if runner.provider != provider {
		t.Error("provider not set correctly")
	}
	if runner.layout != layout {
		t.Errorf("layout = %v, want %v", runner.layout, layout)
	}
}

func TestRunner_Start(t *testing.T) {
	layout := "2006-01-02 15:04:05"
	providerErr := errors.New("provider error")
	fnErr := errors.New("fn error")

	tests := []struct {
		name         string
		provider     *mockProvider
		fn           func(ctx context.Context) error
		wantErr      error
		checkTime    *time.Time
	}{
		{
			name: "empty string from provider runs fn with current time",
			provider: &mockProvider{
				value: "",
			},
			fn: func(ctx context.Context) error {
				return nil
			},
			wantErr: ErrInvalidFaketimeFileContent,
		},
		{
			name: "provider returns error",
			provider: &mockProvider{
				err: providerErr,
			},
			fn: func(ctx context.Context) error {
				t.Error("fn should not be called")
				return nil
			},
			wantErr: providerErr,
		},
		{
			name: "invalid time format returns parse error",
			provider: &mockProvider{
				value: "invalid-time",
			},
			fn: func(ctx context.Context) error {
				t.Error("fn should not be called")
				return nil
			},
			wantErr: ErrInvalidFaketimeFileContent,
		},
		{
			name: "valid time sets fake time and runs fn",
			provider: &mockProvider{
				value: "2024-01-02 15:04:05",
			},
			fn: func(ctx context.Context) error {
				now := time.Now()
				expected := time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)
				if !now.Equal(expected) {
					t.Errorf("time.Now() = %v, want %v", now, expected)
				}
				return nil
			},
			wantErr: nil,
		},
		{
			name: "valid time with @ prefix",
			provider: &mockProvider{
				value: "@2024-03-15 10:30:00",
			},
			fn: func(ctx context.Context) error {
				now := time.Now()
				expected := time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC)
				if !now.Equal(expected) {
					t.Errorf("time.Now() = %v, want %v", now, expected)
				}
				return nil
			},
			wantErr: nil,
		},
		{
			name: "fn returns error is propagated",
			provider: &mockProvider{
				value: "2024-01-02 15:04:05",
			},
			fn: func(ctx context.Context) error {
				return fnErr
			},
			wantErr: fnErr,
		},
		{
			name: "time with ratio option",
			provider: &mockProvider{
				value: "2024-01-02 15:04:05 x2",
			},
			fn: func(ctx context.Context) error {
				now := time.Now()
				expected := time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)
				// With ratio, time should be at or after the start time
				if now.Before(expected) {
					t.Errorf("time.Now() = %v, should be at or after %v", now, expected)
				}
				return nil
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := NewRunner(tt.provider, layout)
			ctx := context.Background()

			err := runner.Start(ctx, tt.fn)

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
		})
	}
}
