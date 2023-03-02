package app

import (
	"fmt"
	"strings"
)

func DiskNameFromStr(code string) string {
	num := Str2Num(code)
	return DiskNameFromNum(num)
}

func DiskNameFromNum(code int) string {
	return fmt.Sprintf("%s%d", DISK_PRE, code)
}

func DiskCodeNumFromName(name string) int {
	code := strings.Replace(name, DISK_PRE, "", 1)
	return Str2Num(code)
}

func DiskCodeStrFromName(name string) string {
	code := DiskCodeNumFromName(name)
	return DiskCodeNum2Str(code)
}

func DiskCodeStr2Num(code string) int {
	return Str2Num(code)
}

func DiskCodeNum2Str(code int) string {
	return fmt.Sprintf("%03d", code)
}
