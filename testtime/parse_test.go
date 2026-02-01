package testtime

import (
	"testing"

	"github.com/akm/time"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name   string
		layout string
		value  string
		want   time.Time
	}{
		{
			name:   "standard layout",
			layout: "2006-01-02 15:04:05",
			value:  "2024-06-15 12:30:45",
			want:   time.Date(2024, 6, 15, 12, 30, 45, 0, time.UTC),
		},
		{
			name:   "date only",
			layout: "2006-01-02",
			value:  "2024-01-02",
			want:   time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			name:   "RFC3339",
			layout: time.RFC3339,
			value:  "2024-06-15T12:30:45Z",
			want:   time.Date(2024, 6, 15, 12, 30, 45, 0, time.UTC),
		},
		{
			name:   "time only",
			layout: "15:04:05",
			value:  "12:30:45",
			want:   time.Date(0, 1, 1, 12, 30, 45, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(t, tt.layout, tt.value)
			assert.Equal(t, tt.want, got)
		})
	}
}
