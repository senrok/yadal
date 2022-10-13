package utils

import (
	"fmt"
	"time"
)

const (
	// RFC 7231#section-7.1.1.1 timetamp format. e.g Tue, 29 Apr 2014 18:30:38 GMT
	rfc822TimeFormat                           = "Mon, 2 Jan 2006 15:04:05 GMT"
	rfc822TimeFormatSingleDigitDay             = "Mon, _2 Jan 2006 15:04:05 GMT"
	rfc822TimeFormatSingleDigitDayTwoDigitYear = "Mon, _2 Jan 06 15:04:05 GMT"
)

func parseTime(t string, formats ...string) (time.Time, error) {
	for _, format := range formats {
		tt, err := time.Parse(format, t)
		if err == nil {
			return tt, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse %s in any of the input formats: %s", t, formats)
}

func ParseRFC7231Time(lastModified string) (time.Time, error) {
	return parseTime(lastModified, rfc822TimeFormat, rfc822TimeFormatSingleDigitDay, rfc822TimeFormatSingleDigitDayTwoDigitYear)
}
