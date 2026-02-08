package time

import (
	std "time"
)

func StdNow() Time {
	return std.Now()
}
