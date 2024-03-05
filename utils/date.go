package utils

import (
	"fmt"
	"os"
	"time"
)

func FormatDate(date string) *time.Time {
	t, err := time.Parse(os.Getenv("DATE_FORMAT"), date)
	if err != nil {
		return nil
	}
	return &t
}

func FormatTime(e string) *time.Time {
	if len(e) == 0 {
		return nil
	}
	t, err := time.Parse(os.Getenv("DATETIME_FORMAT"), e)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &t
}

func GetCurrentTime() time.Time {
	return time.Now()
}

func CompareTimesByGivenMinute(t1 time.Time, storedTime *time.Time, minute int) bool {
	if storedTime == nil {
		return true
	}
	return t1.Sub(*storedTime).Minutes() > float64(minute)
}

func CompareTimesByGivenSecond(t1 time.Time, storedTime *time.Time, second int) bool {
	if storedTime == nil {
		return true
	}
	return t1.Sub(*storedTime).Seconds() > float64(second)
}

func CompareTimesByGivenHour(t1 time.Time, storedTime *time.Time, hour int) bool {
	if storedTime == nil {
		return true
	}
	return t1.Sub(*storedTime).Hours() > float64(hour)
}

func CompareTimesByGivenDay(t1 time.Time, storedTime *time.Time, day int) bool {
	if storedTime == nil {
		return true
	}
	return t1.Sub(*storedTime).Hours() > float64(day*24)
}

func CompareTimesByGivenMonth(t1 time.Time, storedTime *time.Time, month int) bool {
	if storedTime == nil {
		return true
	}
	return t1.Sub(*storedTime).Hours() > float64(month*24*30)
}
