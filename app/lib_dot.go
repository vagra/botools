package app

import (
	"fmt"

	"github.com/qustavo/dotsql"
)

func ReadDotSQL() {
	println("读取 " + DOT_SQL)

	var err error
	g_dot, err = dotsql.LoadFromFile(DOT_SQL)
	Check(err, "读取 %s 失败", DOT_SQL)
}

func DotLatestVersion() int {

	var version int = 1

	for {
		version += 1

		sql_name := DotVersionSQL(version)

		query_map := g_dot.QueryMap()

		_, yes := query_map[sql_name]
		if !yes {
			version -= 1
			return version
		}
	}
}

func DotVersionSQL(version int) string {
	return fmt.Sprintf("%s%d", MIGRATE, version)
}
