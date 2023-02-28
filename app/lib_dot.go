package app

import "github.com/qustavo/dotsql"

func ReadDotSQL() {
	println("读取 " + DOT_SQL)

	var err error
	g_dot, err = dotsql.LoadFromFile(DOT_SQL)
	Check(err, "读取 %s 失败", DOT_SQL)
}
