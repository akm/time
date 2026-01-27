package testtime

import (
	"testing"

	"github.com/akm/time"

	"github.com/stretchr/testify/assert"
)

func TestTravel(t *testing.T) {
	base := time.Now()

	tv := NewTraveler()
	defer tv.Teardown()

	tv.Set(base.Add(2 * time.Hour))

	actual := time.Now()
	assert.Equal(t, base.Add(2*time.Hour), actual)
}
