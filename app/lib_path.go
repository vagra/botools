package app

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}

		Check(err, "检查路径 %s 是否存在时出错", path)
	}

	return true
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}

		Check(err, "检查文件 %s 是否存在时出错", path)
	}

	if info.IsDir() {
		return false
	}

	return true
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}

		Check(err, "检查目录 %s 是否存在时出错", path)
	}

	if !info.IsDir() {
		return false
	}

	return true
}

func IsHidden(path string) bool {

	name := filepath.Base(path)

	if name[0] == '.' {
		return true
	}

	pointer, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false
	}

	attributes, err := syscall.GetFileAttributes(pointer)
	if err != nil {
		return false
	}

	return attributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0
}
