package dateutil

import (
	"github.com/talut/dotenv"
	"time"
)

var defaultDateFormat = "2006-01-02"
var defaultDateTimeFormat = "2006-01-02 15:04:05"

// FormatDate converts string to time.Time using the format from environment variable DATE_FORMAT
func FormatDate(date string) *time.Time {
	t, err := time.Parse(dotenv.GetString("DATE_FORMAT", defaultDateFormat), date)
	if err != nil {
		return nil
	}
	return &t
}

// FormatTime converts string to time.Time using the format from environment variable DATETIME_FORMAT
func FormatTime(e string) *time.Time {
	if len(e) == 0 {
		return nil
	}
	t, err := time.Parse(dotenv.GetString("DATE_TIME_FORMAT", defaultDateTimeFormat), e)
	if err != nil {
		return nil
	}
	return &t
}

// CompareTimesByGivenMinute compares two times by given minute & returns true if t1 is greater than storedTime by given minute
func CompareTimesByGivenMinute(t1 time.Time, storedTime *time.Time, minute int) bool {
	if storedTime == nil {
		return true
	}
	return t1.Sub(*storedTime).Minutes() > float64(minute)
}

// CompareTimesByGivenSecond compares two times by given second & returns true if t1 is greater than storedTime by given second
func CompareTimesByGivenSecond(t1 time.Time, storedTime *time.Time, second int) bool {
	if storedTime == nil {
		return true
	}
	return t1.Sub(*storedTime).Seconds() > float64(second)
}

// CompareTimesByGivenHour compares two times by given hour & returns true if t1 is greater than storedTime by given hour
func CompareTimesByGivenHour(t1 time.Time, storedTime *time.Time, hour int) bool {
	if storedTime == nil {
		return true
	}
	return t1.Sub(*storedTime).Hours() > float64(hour)
}

// CompareTimesByGivenDay compares two times by given day & returns true if t1 is greater than storedTime by given day
func CompareTimesByGivenDay(t1 time.Time, storedTime *time.Time, day int) bool {
	if storedTime == nil {
		return true
	}
	return t1.Sub(*storedTime).Hours() > float64(day*24)
}

// CompareTimesByGivenMonth compares two times by given month & returns true if t1 is greater than storedTime by given month
func CompareTimesByGivenMonth(t1 time.Time, storedTime *time.Time, month int) bool {
	if storedTime == nil {
		return true
	}
	return t1.Sub(*storedTime).Hours() > float64(month*24*30)
}
