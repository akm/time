package testtime

import (
	"testing"

	"github.com/akm/time"

	"github.com/stretchr/testify/assert"
)

func TestSetTimeFunc(t *testing.T) {
	t0 := time.Date(2020, 1, 2, 3, 4, 5, 0, time.FixedLocation)
	f1 := func() time.Time {
		return t0
	}
	defer SetTimeFunc(f1)()

	assert.Equal(t, t0, time.Now())
}

func TestSetTime(t *testing.T) {
	t0 := time.Date(2020, 1, 2, 3, 4, 5, 0, time.FixedLocation)
	t1 := time.Date(2020, 2, 3, 4, 5, 6, 0, time.FixedLocation)

	now := &t0
	defer SetTime(now)()

	assert.Equal(t, t0, time.Now())

	*now = t1
	assert.Equal(t, t1, time.Now())
}
