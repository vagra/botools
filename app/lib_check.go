package app

import (
	"fmt"
	"os"

	"github.com/eiannone/keyboard"
)

func Check(err error, format string, args ...any) {
	if err != nil {
		println(err.Error())
		fmt.Printf(format, args...)
		WaitExit(1)
	}
}

func Exit(code int) {
	println("exit: bye bye.")
	os.Exit(code)
}

func WaitExit(code int) {
	println()
	println("press any key to exit")

	keyboard.GetSingleKey()

	Exit(code)
}
