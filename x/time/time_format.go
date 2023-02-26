package time

import (
	"fmt"
	"time"
)

// ref to "github.com/spf13/cast"

// StringToDate attempts to parse a string into a time.Time type using a
// predefined list of formats.  If no suitable format is found, an error is
// returned.
func StringToDate(s string) (time.Time, error) {
	return parseDateWith(s, time.UTC, timeFormats)
}

// StringToDateInDefaultLocation casts an empty interface to a time.Time,
// interpreting inputs without a timezone to be in the given location,
// or the local timezone if nil.
func StringToDateInDefaultLocation(s string, location *time.Location) (time.Time, error) {
	return parseDateWith(s, location, timeFormats)
}

type timeFormatType int

const (
	timeFormatNoTimezone timeFormatType = iota
	timeFormatNamedTimezone
	timeFormatNumericTimezone
	timeFormatNumericAndNamedTimezone
	timeFormatTimeOnly
)

type timeFormat struct {
	format string
	typ    timeFormatType
}

func (f timeFormat) hasTimezone() bool {
	// We don't include the formats with only named timezones, see
	// https://github.com/golang/go/issues/19694#issuecomment-289103522
	return f.typ >= timeFormatNumericTimezone && f.typ <= timeFormatNumericAndNamedTimezone
}

var (
	timeFormats = []timeFormat{
		{time.RFC3339, timeFormatNumericTimezone},
		{"2006-01-02T15:04:05", timeFormatNoTimezone}, // iso8601 without timezone
		{time.RFC1123Z, timeFormatNumericTimezone},
		{time.RFC1123, timeFormatNamedTimezone},
		{time.RFC822Z, timeFormatNumericTimezone},
		{time.RFC822, timeFormatNamedTimezone},
		{time.RFC850, timeFormatNamedTimezone},
		{"2006-01-02 15:04:05.999999999 -0700 MST", timeFormatNumericAndNamedTimezone}, // Time.String()
		{"2006-01-02T15:04:05-0700", timeFormatNumericTimezone},                        // RFC3339 without timezone hh:mm colon
		{"2006-01-02 15:04:05Z0700", timeFormatNumericTimezone},                        // RFC3339 without T or timezone hh:mm colon
		{"2006-01-02 15:04:05", timeFormatNoTimezone},
		{time.ANSIC, timeFormatNoTimezone},
		{time.UnixDate, timeFormatNamedTimezone},
		{time.RubyDate, timeFormatNumericTimezone},
		{"2006-01-02 15:04:05Z07:00", timeFormatNumericTimezone},
		{"2006-01-02", timeFormatNoTimezone},
		{"02 Jan 2006", timeFormatNoTimezone},
		{"2006-01-02 15:04:05 -07:00", timeFormatNumericTimezone},
		{"2006-01-02 15:04:05 -0700", timeFormatNumericTimezone},
		{time.Kitchen, timeFormatTimeOnly},
		{time.Stamp, timeFormatTimeOnly},
		{time.StampMilli, timeFormatTimeOnly},
		{time.StampMicro, timeFormatTimeOnly},
		{time.StampNano, timeFormatTimeOnly},
	}
)

func parseDateWith(s string, location *time.Location, formats []timeFormat) (d time.Time, e error) {
	for _, format := range formats {
		if d, e = time.Parse(format.format, s); e == nil {
			// Some time formats have a zone name, but no offset, so it gets
			// put in that zone name (not the default one passed in to us), but
			// without that zone's offset. So set the location manually.
			if format.typ <= timeFormatNamedTimezone {
				if location == nil {
					location = time.Local
				}
				year, month, day := d.Date()
				hour, min, sec := d.Clock()
				d = time.Date(year, month, day, hour, min, sec, d.Nanosecond(), location)
			}

			return
		}
	}
	return d, fmt.Errorf("unable to parse date: %s", s)
}

// toInt returns the int value of v if v or v's underlying type
// is an int.
// Note that this will return false for int64 etc. types.
func toInt(v interface{}) (int, bool) {
	switch v := v.(type) {
	case int:
		return v, true
	case time.Weekday:
		return int(v), true
	case time.Month:
		return int(v), true
	default:
		return 0, false
	}
}
