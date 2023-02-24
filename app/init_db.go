package app

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/qustavo/dotsql"
)

func InitDB() error {
	println("start: init db")

	CheckConfig()
	CheckDBsDirExist()

	println()
	ReadSQL()

	println()
	GetNotExistDBs()

	println()
	CreateTables()

	println()
	println("init db done!")
	return nil
}

func CheckConfig() {
	if !ReadConfig() {
		WaitExit(1)
	}
}

func CheckDBsDirExist() {
	if !DirExist(DB_DIR) {
		err := os.Mkdir(DB_DIR, os.ModePerm)
		Check(err, "创建数据库目录 "+DB_DIR+" 时出错")
	}
}

func GetNotExistDBs() {
	g_dbs = make(map[string]*sql.DB)

	for disk_name := range g_disks {

		db_path := GetDBPath(disk_name)

		if DBExist(db_path) {
			fmt.Printf("数据库 %s 已存在，跳过\n", db_path)
			continue
		}

		fmt.Printf("数据库 %s 不存在，新建...\n", db_path)

		db, err := sql.Open("sqlite3", db_path)
		Check(err, "打开数据库 "+db_path+" 失败")

		g_dbs[disk_name] = db
	}
}

func GetDBs() {
	g_dbs = make(map[string]*sql.DB)

	for disk_name := range g_disks {

		db_path := GetDBPath(disk_name)

		fmt.Printf("打开数据库 %s\n", db_path)

		db, err := sql.Open("sqlite3", db_path)
		Check(err, "打开数据库 "+db_path+" 失败")

		g_dbs[disk_name] = db
	}
}

func CreateTables() {
	for db_name, db := range g_dbs {
		fmt.Printf("初始化数据 %s\n", db_name)

		DBCreateDirsTable(db)
		DBCreateFilesTable(db)
		DBCreateInfosTable(db)
	}
}

func ReadSQL() {
	println("读取 " + DOT_SQL)

	var err error
	g_dot, err = dotsql.LoadFromFile(DOT_SQL)
	Check(err, "读取 "+DOT_SQL+" 失败")
}

func DeleteDB() {
	for db_name := range g_dbs {
		db_path := GetDBPath(db_name)

		fmt.Printf("删除数据库 %s", db_path)

		if !FileExist(db_path) {
			println(" (不存在)")
			continue
		}
		println()

		err := os.Remove(db_path)
		Check(err, "删除数据库文件 "+db_path+" 失败")
	}
}

func DBExist(db_path string) bool {
	if FileExist(db_path) {
		return true
	}

	return false
}

func AnyDBExist() bool {
	for db_name := range g_dbs {
		db_path := GetDBPath(db_name)
		if FileExist(db_path) {
			return true
		}
	}

	return false
}

func AllDBExist() bool {
	for db_name := range g_dbs {
		db_path := GetDBPath(db_name)
		if !FileExist(db_path) {
			return false
		}
	}

	return true
}
