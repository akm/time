package testtime

import (
	"github.com/akm/time"
	"github.com/akm/time/internal"
)

func SetTimeFunc(f func() time.Time) func() {
	var backup func() time.Time
	internal.NowFunc, backup = f, internal.NowFunc
	return func() {
		internal.NowFunc = backup
	}
}

func SetTime(v *time.Time) func() {
	return SetTimeFunc(func() time.Time { return *v })
}
