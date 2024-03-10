package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestFormatDate(t *testing.T) {
	os.Unsetenv("DATE_FORMAT") // Ensure DATE_FORMAT is not set
	testDate := "2023-03-10"
	expectedTime, _ := time.Parse(defaultDateFormat, testDate)

	// Test with default date format
	result := FormatDate(testDate)
	assert.NotNil(t, result)
	assert.True(t, expectedTime.Equal(*result))

	// Test with custom date format
	customFormat := "02-01-2006"
	os.Setenv("DATE_FORMAT", customFormat)
	testDateCustomFormat := "10-03-2023"
	expectedTimeCustom, _ := time.Parse(customFormat, testDateCustomFormat)
	resultCustom := FormatDate(testDateCustomFormat)
	assert.NotNil(t, resultCustom)
	assert.True(t, expectedTimeCustom.Equal(*resultCustom))

	// Test with invalid date
	invalidDate := "invalid-date"
	resultInvalid := FormatDate(invalidDate)
	assert.Nil(t, resultInvalid)

	// Cleanup
	os.Unsetenv("DATE_FORMAT")
}

func TestFormatTime(t *testing.T) {
	os.Unsetenv("DATETIME_FORMAT") // Ensure DATETIME_FORMAT is not set
	testDateTime := "2023-03-10 15:04:05"
	expectedTime, _ := time.Parse(defaultDateTimeFormat, testDateTime)

	// Test with default datetime format
	result := FormatTime(testDateTime)
	assert.NotNil(t, result)
	assert.True(t, expectedTime.Equal(*result))

	// Test with empty string
	resultEmpty := FormatTime("")
	assert.Nil(t, resultEmpty)

	// Test with invalid datetime
	invalidDateTime := "invalid-datetime"
	resultInvalid := FormatTime(invalidDateTime)
	assert.Nil(t, resultInvalid)

	// Cleanup
	os.Unsetenv("DATETIME_FORMAT")
}

func TestCompareTimesByGivenMinute(t *testing.T) {
	now := time.Now()
	tenMinutesBefore := now.Add(-10 * time.Minute)
	assert.True(t, CompareTimesByGivenMinute(now, &tenMinutesBefore, 5))
	assert.False(t, CompareTimesByGivenMinute(now, &tenMinutesBefore, 15))
	assert.True(t, CompareTimesByGivenMinute(now, nil, 10))
}

func TestCompareTimesByGivenHour(t *testing.T) {
	now := time.Now()
	threeHoursBefore := now.Add(-3 * time.Hour)

	// Test where t1 is more than 2 hours ahead of storedTime
	assert.True(t, CompareTimesByGivenHour(now, &threeHoursBefore, 2))

	// Test where t1 is not more than 3 hours ahead of storedTime
	assert.False(t, CompareTimesByGivenHour(now, &threeHoursBefore, 3))

	// Test where storedTime is nil
	assert.True(t, CompareTimesByGivenHour(now, nil, 3))
}

func TestCompareTimesByGivenDay(t *testing.T) {
	now := time.Now()
	twoDaysBefore := now.AddDate(0, 0, -2)

	// Test where t1 is more than 1 day ahead of storedTime
	assert.True(t, CompareTimesByGivenDay(now, &twoDaysBefore, 1))

	// Test where t1 is not more than 2 days ahead of storedTime
	assert.False(t, CompareTimesByGivenDay(now, &twoDaysBefore, 2))

	// Test where storedTime is nil
	assert.True(t, CompareTimesByGivenDay(now, nil, 2))
}

func TestCompareTimesByGivenMonth(t *testing.T) {
	now := time.Now()
	oneMonthBefore := now.AddDate(0, -1, 0)

	// Test where t1 is more than 0.5 months ahead of storedTime (approximation)
	assert.True(t, CompareTimesByGivenMonth(now, &oneMonthBefore, 0))

	// Test where t1 is not more than 1 month ahead of storedTime
	assert.False(t, CompareTimesByGivenMonth(now, &oneMonthBefore, 1))

	// Test where storedTime is nil
	assert.True(t, CompareTimesByGivenMonth(now, nil, 1))
}

func TestGetCurrentTime(t *testing.T) {
	now := time.Now()
	assert.True(t, GetCurrentTime().Sub(now).Seconds() < 1)
}
