package util

import "time"

const (
	defaultFormat = "2006-01-02 15:04:05"
)

func TimeToString(time time.Time) string {
	return time.Format(defaultFormat)
}
