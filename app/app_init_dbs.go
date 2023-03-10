package app

import (
	"fmt"
)

func InitDBs() error {
	println("start: init db")

	println()
	ReadConfig()

	println()
	ReadDotSQL()

	println()
	CheckDBsDirExists()

	println()
	GetNotInitedDBs()

	CheckTaskHasDBs()

	println()
	STInitDBs()

	println()
	println("init db done!")
	return nil
}

func STInitDBs() {
	println("初始化数据库")

	for db_name, db := range g_dbs {
		db_path := GetDBPath(db_name)
		fmt.Printf("在数据库 %s 中创建表\n", db_path)

		InitDBWorker(db)
	}
}

func InitDBWorker(db *DB) {
	db.CreateDirsTable()
	db.CreateFilesTable()
	db.CreateVDirsTable()
	db.CreateVFilesTable()
	db.CreateInfosTable()
}
