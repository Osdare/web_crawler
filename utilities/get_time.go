package utilities

import "time"

func GetTimeInt() int64 {
	return time.Now().Unix()
}