package timeutil

import (
	"strconv"
	"time"
)

var (
	langDurationFormats = map[string][]string{
		"zh-CN": {"日", "小时", "分钟", "秒"},
		"en-US": {"d", "h", "m", "s"},
	}

	defaultLang = "zh-CN"
)

func FormatDuration(d time.Duration) string {
	return FormatDurationWithLang(d, defaultLang)
}

func FormatDurationWithLang(d time.Duration, lang string) string {
	formats, ok := langDurationFormats[lang]
	if !ok {
		formats = langDurationFormats[defaultLang]
	}

	tsec := int(d.Seconds())

	seconds := tsec % 60
	minutes := tsec / 60

	hours := minutes / 60
	minutes = minutes % 60

	day := hours / 24
	hours = hours % 24

	formatStr := ""
	if day != 0 {
		formatStr += strconv.Itoa(day) + formats[0]
	}

	if hours != 0 {
		formatStr += strconv.Itoa(hours) + formats[1]
	}

	if minutes != 0 {
		formatStr += strconv.Itoa(minutes) + formats[2]
	}

	if len(formatStr) == 0 || seconds != 0 {
		formatStr += strconv.Itoa(seconds) + formats[3]
	}
	return formatStr
}
