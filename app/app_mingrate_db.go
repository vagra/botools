package app

import (
	"database/sql"
	"fmt"
)

func MigrateDB() error {
	println("start: migrate db")

	println()
	ReadConfig()

	println()
	ReadDotSQL()

	println()
	CheckDBsDirExists()

	println()
	CheckAllDBExists()

	println()
	CheckAllDBInited()

	println()
	GetInitedDBs()

	CheckTaskHasDBs()

	println()
	ConfirmMigrateDB()

	println()
	STMigrateDB()

	println()
	println("migrate db done!")
	return nil
}

func STMigrateDB() {

	g_latest = DotLatestVersion()
	fmt.Printf("数据库最新版本 v%d\n", g_latest)
	println()

	for name := range g_dbs {
		MigrateDBWorker(name)
		println()
	}
}

func MigrateDBWorker(disk_name string) {

	var db *DB = g_dbs[disk_name]

	db_path := GetDBPath(disk_name)

	fmt.Printf("%s worker: start migrate db %s\n", disk_name, db_path)

	if !db.InfosTableExists() {
		db.CreateInfosTable()
	}

	old_ver := db.GetVersion()
	fmt.Printf("数据库当前版本 v%d\n", old_ver)

	new_ver := old_ver

	for {
		new_ver += 1
		if new_ver > g_latest {
			new_ver -= 1
			break
		}

		sql_name := DotVersionSQL(new_ver)
		fmt.Printf("执行数据库升级命令 %s\n", sql_name)

		_, err := g_dot.Exec((*sql.DB)(db), sql_name)
		Check(err, "执行数据库升级命令 %s 时失败", sql_name)
	}

	if new_ver > old_ver {
		fmt.Printf("数据库版本更新到 v%d\n", new_ver)
	} else {
		println("数据库没有更新")
	}

	fmt.Printf("%s worker: stop.\n", disk_name)
}

func ConfirmMigrateDB() {
	println("本程序用于升级现有数据库到最新版本")
	println("您确定要执行这个操作吗？请输入 yes 或 no ：")

	CheckConfirm()
}
