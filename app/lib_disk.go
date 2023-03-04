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

// disk-008 -> 008
func DiskCodeFromName(disk_name string) string {
	return strings.Replace(disk_name, DISK_PRE, "", 1)
}

// 008 -> disk-008
func DiskNameFromCode(code string) string {
	return fmt.Sprintf("%s%s", DISK_PRE, code)
}
