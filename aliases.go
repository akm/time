package time

import (
	// https://pkg.go.dev/time
	orig "time"
)

type (
	Duration   = orig.Duration
	Location   = orig.Location
	ParseError = orig.ParseError
	Month      = orig.Month
	Ticker     = orig.Ticker
	Time       = orig.Time
	Timer      = orig.Timer
	Weekdday   = orig.Weekday
)

var (
	UTC   = orig.UTC
	Local = orig.Local
)

var (
	After = orig.After
	Sleep = orig.Sleep
	Tick  = orig.Tick

	// Duration
	ParseDuration = orig.ParseDuration
	Since         = orig.Since
	Until         = orig.Until

	// Location
	FixedZone              = orig.FixedZone
	LoadLocation           = orig.LoadLocation
	LoadLocationFromTZData = orig.LoadLocationFromTZData

	// Ticker
	NewTicker = orig.NewTicker

	// Time
	Date = orig.Date
	// Now = orig.Now // replaced in ./now.go
	Parse           = orig.Parse
	ParseInLocation = orig.ParseInLocation
	Unix            = orig.Unix
	UnixMicro       = orig.UnixMicro
	UnixMilli       = orig.UnixMilli

	// Timer
	AfterFunc = orig.AfterFunc
	NewTimer  = orig.NewTimer
)

const (
	Layout      = orig.Layout
	ANSIC       = orig.ANSIC
	UnixDate    = orig.UnixDate
	RubyDate    = orig.RubyDate
	RFC822      = orig.RFC822
	RFC822Z     = orig.RFC822Z
	RFC850      = orig.RFC850
	RFC1123     = orig.RFC1123
	RFC1123Z    = orig.RFC1123Z
	RFC3339     = orig.RFC3339
	RFC3339Nano = orig.RFC3339Nano
	Kitchen     = orig.Kitchen
	Stamp       = orig.Stamp
	StampMilli  = orig.StampMilli
	StampMicro  = orig.StampMicro
	StampNano   = orig.StampNano
	DateTime    = orig.DateTime
	DateOnly    = orig.DateOnly
	TimeOnly    = orig.TimeOnly

	Nanosecond  = orig.Nanosecond
	Microsecond = orig.Microsecond
	Millisecond = orig.Millisecond
	Second      = orig.Second
	Minute      = orig.Minute
	Hour        = orig.Hour

	// Month
	January   = orig.January
	February  = orig.February
	March     = orig.March
	April     = orig.April
	May       = orig.May
	June      = orig.June
	July      = orig.July
	August    = orig.August
	September = orig.September
	October   = orig.October
	November  = orig.November
	December  = orig.December
)
