package app

import (
	"fmt"
	"log"
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

func PassErr(err error, format string, args ...any) bool {
	ok := err == nil

	return PassOk(ok, format, args...)
}

func PassOk(ok bool, format string, args ...any) bool {
	if !ok {
		log.Printf(format, args...)
		log.Println()
		fmt.Printf(format, args...)
		fmt.Println()
	}

	return ok
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
