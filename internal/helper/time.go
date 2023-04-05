package helper

import (
	"strconv"
	"time"
)

// ParseTime is used to converting a time of string type
func ParseTime(t string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", t, time.Local)
}

// ParseDate is used to converting a time of string type
func ParseDate(t string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", t, time.Local)
}

// SetTimeMidNight is uset to setting time 00:00:00
// ex) t: 2022-10-28 19:30:33 -> return 2022-10-28 00:00:00
func SetTimeMidNight(t time.Time) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", t.Format("2006-01-02"), time.Local)
}

// IsValidDate is used to checking start & end date's validation.
// ex) start: 2022-10-28 end: 2022-10-28 or 2022-10-29 -> true
//
//	end: 2022-10-29 end: 2022-10-28 -> false
func IsValidDate(start time.Time, end time.Time) bool {
	if res := start.Sub(end); res.Hours() > 0 {
		return false
	}
	return true
}

// IsToday is used to checking time is same with today date.
// ex) today date: 2022-10-28 parameter: 2022-10-28 -> true
//
//	today date: 2022-10-28 parameter: 2022-10-27 -> false
func IsToday(t time.Time) bool {
	now := time.Now()
	if now.Day() != t.Day() ||
		now.Month() != t.Month() ||
		now.Year() != t.Year() {
		return false
	}
	return true
}

// In case the column of the database table is 'datetime', the timezone of the time is changed to UTC.
func ToUTCDatetime(org time.Time) time.Time {
	_, offset := org.Zone()

	return org.In(time.UTC).Add(time.Duration(offset) * time.Second)
}

// 입력 시간이 몇년 몇월에 몇주차인지 제공
func NumberOfWeekInMonth(now time.Time) string {
	beginningOfTheMonth := time.Date(now.Year(), now.Month(), 1, 1, 1, 1, 1, time.UTC)
	thisYear, thisWeek := now.ISOWeek()
	beginYear, beginningWeek := beginningOfTheMonth.ISOWeek()

	ret := 1 + thisWeek - beginningWeek

	if thisYear != beginYear {
		ret = thisWeek
	}
	return now.Format("2006-01 ") + strconv.Itoa(ret)
}
