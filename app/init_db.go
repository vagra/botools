package app

import (
	"database/sql"
	"os"

	"github.com/qustavo/dotsql"
)

func InitDB() error {
	println("start: init db")

	if !ReadConfig() {
		WaitExit(1)
	}

	CheckDBsDir()

	println()
	GetDBs()
	ReadSQL()

	if AnyDBExist() {
		println()
		println("初始化数据库会删除现有数据库，请谨慎操作！")
		println("您确定要删除现有的数据库文件？请输入 yes 或 no ：")

		if Confirm() {
			println()
			DeleteDB()
		} else {
			WaitExit(1)
		}
	}

	println()
	CreateTables()

	println()
	println("init db done!")
	return nil
}

func CheckDBsDir() {
	if !DirExist(DB_DIR) {
		err := os.Mkdir(DB_DIR, os.ModePerm)
		Check(err, "创建数据库目录 "+DB_DIR+" 时出错")
	}
}

func GetDBs() {
	g_dbs = make(map[string]*sql.DB)

	for disk_name := range g_disks {

		db_name := GetDBName(disk_name)

		println("打开数据库 " + db_name)

		db, err := sql.Open("sqlite3", db_name)
		Check(err, "打开数据库 "+db_name+" 失败")

		g_dbs[db_name] = db
	}
}

func CreateTables() {
	for db_name, db := range g_dbs {
		println("初始化数据 " + db_name)

		_, err := g_dot.Exec(db, SQL_CREATE_DIRS)
		Check(err, "在数据库 "+db_name+" 中创建 dirs 表失败")

		_, err = g_dot.Exec(db, SQL_CREATE_FILES)
		Check(err, "在数据库 "+db_name+" 中创建 files 表失败")

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
		print("删除数据库 " + db_name)

		if !FileExist(db_name) {
			println(" (不存在)")
			continue
		}
		println()

		err := os.Remove(db_name)
		Check(err, "删除数据库文件 "+db_name+" 失败")
	}
}

func AnyDBExist() bool {
	for db_name := range g_dbs {
		if FileExist(db_name) {
			return true
		}
	}

	return false
}

func AllDBExist() bool {
	for db_name := range g_dbs {
		if !FileExist(db_name) {
			return false
		}
	}

	return true
}
