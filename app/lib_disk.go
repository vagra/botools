package app

import (
	"fmt"
	"regexp"
	"strings"
)

func CheckHasDisksConfig() {
	if !HasDisksConfig() {
		fmt.Printf("请在 %s 中设置 disks 列表，格式：\n", CONFIG_INI)
		println("[disks]")
		println("disk-name = disk:/path/")
		println("disk-name = disk:/path/")
		println("disk-name 必须以 'disk-' 开头，后跟若干位英文字母或数字")
		println("disk-path 必须是以盘符开头的绝对路径，并且要使用正斜杠 /")
		WaitExit(1)
	}
}

func CheckDiskNameValid(name string) {
	if !DiskNameValid(name) {
		println("disk-name 必须以 'disk-' 开头，后跟若干位英文字母或数字")
		fmt.Printf("请检查 %s\n", CONFIG_INI)
		WaitExit(1)
	}
}

func CheckDiskPathValid(path string) {
	if !DiskPathValid(path) {
		println("disk-path 必须是以盘符开头的绝对路径，并且要使用正斜杠 /")
		fmt.Printf("请检查 %s\n", CONFIG_INI)
		WaitExit(1)
	}
}

func HasDisksConfig() bool {
	return len(g_disks) > 0
}

func DiskNameValid(name string) bool {
	if len(name) <= 0 {
		return false
	}

	return regexp.MustCompile(`^disk-[a-zA-Z0-9]+$`).MatchString(name)
}

func DiskPathValid(path string) bool {
	if len(path) <= 0 {
		return false
	}

	if strings.Contains(path, "\\") {
		return false
	}

	if !strings.Contains(path, ":") {
		return false
	}

	return true
}

// disk-001 -> disk-1
func DiskNameStr2Num(name string) string {
	num := DiskCodeNumFromName(name)
	return DiskNameFromNum(num)
}

// disk-1 -> disk-001
func DiskNameNum2Str(name string) string {
	code := DiskCodeStrFromName(name)
	return DiskNameFromStr(code)
}

// 001 -> disk-1
func DiskNameFromStr(code string) string {
	num := Str2Num(code)
	return DiskNameFromNum(num)
}

// 1 -> disk-1
func DiskNameFromNum(code int) string {
	return fmt.Sprintf("%s%d", DISK_PRE, code)
}

// disk-008 -> 8
func DiskCodeNumFromName(name string) int {
	code := strings.Replace(name, DISK_PRE, "", 1)
	return Str2Num(code)
}

// disk-008 -> 008
func DiskCodeStrFromName(name string) string {
	code := DiskCodeNumFromName(name)
	return DiskCodeNum2Str(code)
}

// 006 -> 6
func DiskCodeStr2Num(code string) int {
	return Str2Num(code)
}

// 6 -> 006
func DiskCodeNum2Str(code int) string {
	return fmt.Sprintf("%03d", code)
}
