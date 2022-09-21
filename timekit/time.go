package timekit

import (
	"fmt"
	"time"
)

const (
	SimpleDateTime = "2006-01-02T15:04:05"
	SimpleDate     = "2006-01-02"
	SimpleTime     = "15:04:05"
)

var timeFormats = []string{
	"20060102",
	"200601021504",
	"20060102150405",
	"2006-01-02 15:04",
	"2006-01-02 15:04:05",
	time.RFC3339,
	"2006-01-02T15:04:05", // iso8601 without timezone
	time.RFC1123Z,
	time.RFC1123,
	time.RFC822Z,
	time.RFC822,
	time.RFC850,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	"2006-01-02 15:04:05.999999999 -0700 MST", // Time.String()
	"2006-01-02",
	"02 Jan 2006",
	"2006-01-02T15:04:05-0700", // RFC3339 without timezone hh:mm colon
	"2006-01-02 15:04:05 -07:00",
	"2006-01-02 15:04:05 -0700",
	"2006-01-02 15:04:05Z07:00", // RFC3339 without T
	"2006-01-02 15:04:05Z0700",  // RFC3339 without T or timezone hh:mm colon
	"2006-01-02 15:04:05.000",
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
}
var timeFormatsForExpandSection = [][]string{
	// 按年
	{"2006", "2006年"},
	// 按月
	{"200601", "2006-01", "Jan 2006", "2006年01月"},
	// 按天
	{"20060102", "2006-01-02", "02 Jan 2006", "2006年01月02日"},
	// 时
	{"2006010215", "2006-01-02 15", "02 Jan 2006", "2006年01月02日"},
	// 分
	{"200601021504", "2006-01-02 15:04", "2006-01-02T15:04"},
	// 秒，其他
	{
		"20060102150405",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05", // iso8601 without timezone
		time.RFC3339,
		time.RFC1123Z,
		time.RFC1123,
		time.ANSIC,
		time.UnixDate,
		"2006-01-02 15:04:05.999999999 -0700 MST", // Time.String()
		"2006-01-02T15:04:05-0700",                // RFC3339 without timezone hh:mm colon
		"2006-01-02 15:04:05 -07:00",
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05Z07:00", // RFC3339 without T
		"2006-01-02 15:04:05Z0700",  // RFC3339 without T or timezone hh:mm colon
		"2006-01-02 15:04:05.000",
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	},
}

// StringToDate attempts to parse a string into a time.Time type using a
// predefined list of formats.  If no suitable format is found, an error is
// returned.
func StringToDate(s string) (time.Time, error) {
	return parseDateWith(s, timeFormats)
}

func parseDateWith(s string, dates []string) (d time.Time, e error) {
	for _, dateType := range dates {
		if d, e = time.ParseInLocation(dateType, s, time.Local); e == nil {
			return
		}
	}
	return d, fmt.Errorf("unable to parse date: %s", s)
}

// Str2Time 字符串转成时间（或日期），如果 无效，则返回 0时间  d.IsZero()
func Str2Time(s string) (d time.Time) {
	var e error
	if s == "" {
		return d
	}
	for _, dateType := range timeFormats {
		if d, e = time.ParseInLocation(dateType, s, time.Local); e == nil {
			return
		}
	}
	return d
}

// Str2TimeExpand 扩展包含当前的时段；如参数是按天，则扩展到明天；如果参数精确到小时，则扩展到下一小时
func Str2TimeExpand(s string) (d time.Time) {
	for i := range timeFormatsForExpandSection {
		for _, dateType := range timeFormatsForExpandSection[i] {
			if t1, e1 := time.ParseInLocation(dateType, s, time.Local); e1 == nil {
				year, month, day := t1.Date()
				hour, min, sec := t1.Clock()
				switch i {
				case 0:
					d = time.Date(year+1, 0, 0, 0, 0, 0, 0, t1.Location())
				case 1:
					d = time.Date(year, month+1, 0, 0, 0, 0, 0, t1.Location())
				case 2:
					d = time.Date(year, month, day+1, 0, 0, 0, 0, t1.Location())
				case 3:
					d = time.Date(year, month, day, hour+1, 0, 0, 0, t1.Location())
				case 4:
					d = time.Date(year, month, day, hour, min+1, 0, 0, t1.Location())
				case 5:
					d = time.Date(year, month, day, hour, min, sec+1, 0, t1.Location())
				default:
					d = t1
				}
				return d
			}
		}
	}
	return Str2Time(s)
}

// Str2Date 字符串转成日期（时、分、秒为0），如果无效，则返回 0时间  d.IsZero()
func Str2Date(s string) (d time.Time) {
	if s == "" {
		return
	}
	d, e := parseDateWith(s, timeFormats)
	if e == nil && !d.IsZero() {
		y, m, day := d.Date()
		d = time.Date(y, m, day, 0, 0, 0, 0, d.Location())
	}
	return
}

// Truncate4Day 按天截断
func Truncate4Day(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// TimeTruncate 截断对齐；不支持按月、按年
//  如果是按其他维度，比如 一天内小时、小时内分钟 等，则要先对齐到目标点，再使用差值计算。
// 用法:
//  t1 => 2021-06-17 18:32:11 +0800 CST
//  按天: TimeTruncate(t1, time.Hour*24) => 2021-06-17 00:00:00 +0800 CST
//  按7小时: TimeTruncate(t1, time.Hour*7) => 2021-06-17 12:00:00 +0800 CST
//  按7小时（当天）: TimeTruncate2(t1, time.Hour*24, time.Hour*7) => 2021-06-17 14:00:00 +0800 CST
//  按小时: TimeTruncate(t1, time.Hour) => 2021-06-17 18:00:00 +0800 CST
//  按分钟: TimeTruncate(t1, time.Minute) => 2021-06-17 18:32:00 +0800 CST
func TimeTruncate(t0 time.Time, d time.Duration) time.Time {
	_, dif := t0.Zone()
	addDiff := time.Second * time.Duration(dif)
	return t0.Add(addDiff).Truncate(d).Add(-addDiff)
}

// 依次对齐
func TimeTruncate2(t0 time.Time, d1, d2 time.Duration) time.Time {
	_, dif := t0.Zone()
	addDiff := time.Second * time.Duration(dif)
	t1 := t0.Add(addDiff).Truncate(d1).Add(-addDiff)
	return t1.Add(t0.Sub(t1).Truncate(d2))
}

// TimeMonthStart 取月初日期
func TimeMonthStart(t0 time.Time) time.Time {
	if t0.IsZero() {
		return t0
	}
	y, m, _ := t0.Date()
	return time.Date(y, m, 0, 0, 0, 0, 0, t0.Location())
}

func TimeMonthAdd(t0 time.Time, inc int) time.Time {
	if t0.IsZero() {
		return t0
	}
	y, m, d := t0.Date()
	return time.Date(y, m+time.Month(inc), d, 0, 0, 0, 0, t0.Location())
}
