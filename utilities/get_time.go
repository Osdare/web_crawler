package utilities

import "time"

func GetTimeInt() int {
	return int(time.Now().Unix())
}