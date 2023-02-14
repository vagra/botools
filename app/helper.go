package app

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

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

func Confirm() bool {

	var input string
	var yes string

	for {
		_, err := fmt.Scanln(&input)
		if err != nil {
			println("请输入")
			continue
		}

		yes = strings.ToLower(input)

		switch yes {
		case "yes":
			return true
		case "n":
			fallthrough
		case "no":
			return false
		default:
			println("请输入 yes 或 no")
		}
	}
}

func GenUID(prefix string, counter *int64) string {
	*counter += 1
	return fmt.Sprintf("%s-%08d", prefix, *counter)
}

func GetUID(prefix string, id int64) string {
	return fmt.Sprintf("%s-%016d", prefix, id)
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

func IsValidName(name string) bool {
	return regexp.MustCompile(`^[a-zA-Z]+[a-zA-Z0-9_\-]+$`).MatchString(name)
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

func GetDBPath(disk_name string) string {
	return fmt.Sprintf("%s/%s%s", DB_DIR, disk_name, DB_EXT)
}

func GetSHA1(path string) (string, int8) {

	file, err := os.Open(path)
	if err != nil {
		return "", 1
	}

	sha1h := sha1.New()
	io.Copy(sha1h, file)
	sum := hex.EncodeToString(sha1h.Sum(nil))

	return sum, 0
}
