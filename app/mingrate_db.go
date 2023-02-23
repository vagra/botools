package app

import (
	"database/sql"
	"fmt"
)

func MigrateDB() error {
	println("start: migrate db")

	CheckConfig()

	println()
	GetDBs()
	ReadSQL()

	println()
	CheckAllDBExist()

	if !ConfirmMigrateDB() {
		WaitExit(1)
	}

	STMigrateDB()

	println()
	println("migrate db done!")

	return nil
}

func STMigrateDB() {

	g_latest = DotLatestVersion()
	fmt.Printf("数据库最新版本 v%d\n", g_latest)

	for name, path := range g_disks {
		print(path)

		MigrateDBWorker(name, path)
	}
}

func MigrateDBWorker(disk_name string, disk_path string) {

	var db *sql.DB = g_dbs[disk_name]

	fmt.Printf("%s worker: start migrate db\n", disk_name)

	if !DBInfosTableExists(db) {
		DBCreateInfosTable(db)
		DBAddInfo(db, 1)
	}

	old_ver := DBGetVersion(db)
	fmt.Printf("数据库当前版本 v%d\n", old_ver)

	new_ver := old_ver

	for {
		new_ver += 1
		if new_ver > g_latest {
			new_ver -= 1
			break
		}

		sql_name := GetVersionSQL(new_ver)
		fmt.Printf("执行数据库升级命令 %s\n", sql_name)

		_, err := g_dot.Exec(db, sql_name)
		Check(err, "执行数据库升级命令 "+sql_name+" 时失败")
	}

	if new_ver > old_ver {
		DBUpdateVersion(db, new_ver)
		fmt.Printf("数据库版本更新到 v%d\n", new_ver)
	} else {
		println("数据库没有更新")
	}

	fmt.Printf("%s worker: stop.\n", disk_name)
}

func ConfirmMigrateDB() bool {
	println()
	println("本程序用于升级现有数据库到最新版本\n")
	println("您确定要执行这个操作吗？请输入 yes 或 no ：")

	return Confirm()
}

func DotLatestVersion() int {

	var version int = 1

	for {
		version += 1

		sql_name := GetVersionSQL(version)
		println(sql_name)

		query_map := g_dot.QueryMap()

		stmt, ok := query_map[sql_name]
		if !ok {
			version -= 1
			return version
		}

		println(stmt)
	}
}

func GetVersionSQL(version int) string {
	return fmt.Sprintf("%s%d", MIGRATE, version)
}