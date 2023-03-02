package app

import (
	"fmt"
)

func DiskNameFromStr(code string) string {
	return DiskNameFromNum(Str2Num(code))
}

func DiskNameFromNum(code int) string {
	return fmt.Sprintf("%s%d", DISK_PRE, code)
}

func DiskCodeStr2Num(code string) int {
	return Str2Num(code)
}

func DiskCodeNum2Str(code int) string {
	return fmt.Sprintf("%03d", code)
}
