package main

import (
	"botools/app"

	_ "github.com/eiannone/keyboard"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/qustavo/dotsql"
	_ "gopkg.in/ini.v1"
)

func main() {
	app.Welcome()
}
