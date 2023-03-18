package app

import (
	"fmt"
	"strconv"
	"strings"
)

func CheckConfirm() {
	if !Confirm() {
		WaitExit(0)
	}
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

func Str2Num(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}

	return num
}

func Num2Str(num int) string {
	return fmt.Sprintf("%d", num)
}

func Str2Int64(str string) int64 {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}

	return num
}

func Int642Str(num int64) string {
	return fmt.Sprintf("%d", num)
}
