package app

import (
	"database/sql"
	"fmt"
)

func InitDBs() error {
	println("start: init db")

	println()
	ReadConfig()
	CheckDBsDirExist()

	println()
	ReadDotSQL()

	println()
	GetNotExistDBs()

	println()
	STInitDBs()

	println()
	println("init db done!")
	return nil
}

func STInitDBs() {
	if len(g_dbs) <= 0 {
		println("没有需要初始化的数据库")
		return
	}

	println("初始化数据库")

	for db_name, db := range g_dbs {
		db_path := GetDBPath(db_name)
		fmt.Printf("在数据库 %s 中创建表\n", db_path)

		InitDBWorker(db)
	}
}

func InitDBWorker(db *sql.DB) {
	DBCreateDirsTable(db)
	DBCreateFilesTable(db)
	DBCreateInfosTable(db)
}
