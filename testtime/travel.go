package testtime

import "github.com/akm/time"

type Traveler struct {
	now      *time.Time
	teardown func()
}

func NewTraveler() *Traveler {
	now := time.Now()
	return &Traveler{now: &now, teardown: SetTime(&now)}
}

func (tv *Traveler) Teardown() {
	tv.teardown()
}

func (tv *Traveler) Set(v time.Time) {
	*tv.now = v
}
