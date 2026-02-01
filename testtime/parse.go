package testtime

import (
	"testing"

	"github.com/akm/time"
)

func Parse(t *testing.T, layout, value string) time.Time {
	tt, err := time.Parse(layout, value)
	if err != nil {
		t.Fatalf("Parse(%q, %q) error: %v", layout, value, err)
	}
	return tt
}
