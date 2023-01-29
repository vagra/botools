package app

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/qustavo/dotsql"
)

func InitDB() error {
	println("start: init db")

	if AnyDBExist() {
		if Confirm() {
			DeleteDB()
		} else {
			WaitExit(1)
		}
	}

	ReadSQL()
	GetDBs()
	CreateTables()

	println("init db done!")
	return nil
}

func GetDBs() {
	var err error

	for i := 0; i < DB_COUNT; i++ {
		println("打开数据库 " + g_db_names[i])

		g_dbs[i], err = sql.Open("sqlite3", g_db_names[i])
		Check(err, "打开数据库 "+g_db_names[i]+" 失败")
	}

	g_disks_db = g_dbs[0]
	g_files_db = g_dbs[1]
	g_dirs_db = g_dbs[2]
	g_file_metas_db = g_dbs[3]
	g_dir_metas_db = g_dbs[4]
}

func CreateTables() {
	var err error
	for i := 0; i < DB_COUNT; i++ {
		println("在数据库 " + g_db_names[i] + " 中创建表 " + g_db_tables[i])

		_, err = g_dot.Exec(g_dbs[i], g_create_sqls[i])
		Check(err, "在数据库 "+g_db_names[i]+" 中创建表 "+g_db_tables[i]+" 失败")
	}
}

func ReadSQL() {
	println("读取 " + DOT_SQL)

	var err error
	g_dot, err = dotsql.LoadFromFile(DOT_SQL)
	Check(err, "读取 "+DOT_SQL+" 失败")
}

func DeleteDB() {
	for _, path := range g_db_names {
		print("删除数据库 " + path)

		if !FileExist(path) {
			println(" (不存在)")
			continue
		}
		println()

		err := os.Remove(path)
		Check(err, "删除数据库文件 "+path+" 失败")
	}
}

func Confirm() bool {
	println("初始化数据库会删除现有数据库，请谨慎操作！")
	println("您确定要删除现有的数据库文件？请输入 yes 或 no ：")

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

func AnyDBExist() bool {
	for _, path := range g_db_names {
		if FileExist(path) {
			return true
		}
	}

	return false
}

func AllDBExist() bool {
	for _, path := range g_db_names {
		if !FileExist(path) {
			return false
		}
	}

	return true
}
