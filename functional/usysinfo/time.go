package usysinfo

import (
	"time"
)

func FormatTime(date time.Time, layout ...string) string {
	timeFormat := "2006/01/02 15:04:05"
	if len(layout) > 0 {
		timeFormat = layout[0]
	}
	cusTime := date.In(getTimeZone())
	return cusTime.Format(timeFormat)
}

func StrToTime(str, layout string) (time.Time, error) {
	return time.ParseInLocation(layout, str, getTimeZone())
}

func getTimeZone() *time.Location {
	return time.FixedZone("UTC+8", 8*60*60)
}
