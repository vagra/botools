package app

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func PassMakeParentDirs(path string) bool {
	ok := MakeParentDirs(path)
	return PassOk(ok, "为 %s 创建上级目录失败", path)
}

func MakeParentDirs(path string) bool {
	path = CleanPath(path)
	names := strings.Split(path, "/")

	if len(names) <= 1 {
		return false
	}

	var parent_path string = names[0]

	for i := 1; i < len(names)-1; i++ {
		parent_path += "/" + names[i]
		parent_path = CleanPath(parent_path)
		if !PassMakeDir(parent_path) {
			return false
		}
	}

	return true
}

func PassMakeDir(path string) bool {
	if DirExists(path) {
		return true
	}

	err := os.Mkdir(path, os.ModePerm)
	return PassErr(err, "创建文件夹 %s 时出错", path)
}

func PassFileExists(path string) bool {
	ok := FileExists(path)
	return PassOk(ok, "文件 %s 不存在", path)
}

func PassCopyFile(src string, dst string) bool {
	return CopyFile(src, dst)
}

func CopyFile(src string, dst string) bool {
	src_file, err := os.Open(src)
	if !PassErr(err, "copy: 打开 src 文件时出错 %s", src) {
		return false
	}
	defer src_file.Close()

	dst_file, err := os.Create(dst)
	if !PassErr(err, "copy: 创建 dst 文件时出错 %s", dst) {
		return false
	}
	defer dst_file.Close()

	_, err = io.Copy(dst_file, src_file)
	if !PassErr(err, "copy: 复制文件时出错 %s", dst) {
		return false
	}

	err = dst_file.Sync()
	if !PassErr(err, "copy: 复制文件时出错 %s", dst) {
		return false
	}

	return true
}

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
