package app

import (
	"errors"
	"os"

	"github.com/eiannone/keyboard"
)

func Check(err error, msg string) {
	if err != nil {
		println(err.Error())
		println(msg)
		WaitExit(1)
	}
}

func Exit(code int) {
	println("exit: bye bye.")
	os.Exit(code)
}

func WaitExit(code int) {
	println("按任意键退出")

	_, _, err := keyboard.GetSingleKey()
	if err != nil {
		panic(err)
	}

	Exit(code)
}

func PathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}

		Check(err, "检查路径 "+path+" 是否存在时出错")
	}

	return true
}

func FileExist(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}

		Check(err, "检查文件 "+path+" 是否存在时出错")
	}

	if info.IsDir() {
		return false
	}

	return true
}

func DirExist(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}

		Check(err, "检查目录 "+path+" 是否存在时出错")
	}

	if !info.IsDir() {
		return false
	}

	return true
}
