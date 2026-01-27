package time

import (
	"github.com/akm/time/internal"
)

func Now() Time {
	return NowWithLocation()
}

func NowWithoutLocation() Time {
	return internal.NowFunc()
}

var FixedLocation = func() *Location {
	return FixedZone("Asia/Tokyo", 9*60*60)
}()

func NowWithLocation() Time {
	return internal.NowFunc().In(FixedLocation)
}
