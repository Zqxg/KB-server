package utils

import (
	"time"
)

// 常见时间格式类型常量
const (
	FormatDate           = "2006-01-02"                // yyyy-mm-dd
	FormatTime           = "15:04:05"                  // hh:mm:ss
	FormatDateTime       = "2006-01-02 15:04:05"       // yyyy-mm-dd hh:mm:ss
	FormatDateTimeWithTZ = "2006-01-02 15:04:05-07:00" // 带时区 yyyy-mm-dd hh:mm:ss-07:00
)

// TimeFormat 通过传入的格式对 time.Time 类型的时间进行格式化
func TimeFormat(t time.Time, format string) string {
	return t.Format(format)
}

// ParseTime 根据指定格式解析时间字符串
func ParseTime(timeStr, format string) (time.Time, error) {
	return time.Parse(format, timeStr)
}
