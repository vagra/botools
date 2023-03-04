package app

import (
	"database/sql"
	"fmt"
	"sync"
)

func TrimIDs() error {
	println("start: trim dir and file ids")

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
	CheckAllDBHasData()

	println()
	GetHasDataDBs()

	CheckTaskHasDBs()

	println()
	ConfirmModDiskIDs()

	println()
	MTModDiskIDs()

	println()
	println("trim dir and file ids done!")
	return nil
}

func MTTrimIDs() {
	println("每个 disk 启动一个线程，截短 dirs 和 files 的 id 和 parent_id 到 8 位")

	var wg sync.WaitGroup

	for name := range g_dbs {
		wg.Add(1)
		go ModDiskIDsWorker(&wg, name)
	}

	wg.Wait()
}

func TrimIDsWorker(wg *sync.WaitGroup, disk_name string) {
	defer wg.Done()

	var db *sql.DB = g_dbs[disk_name]

	db_path := GetDBPath(disk_name)

	fmt.Printf("%s worker: start trim ids %s\n", disk_name, db_path)

	DBTrimDirsID(db)
	DBTrimFilesID(db)

	fmt.Printf("%s worker: stop.\n", disk_name)
}

func ConfirmTrimIDs() {
	println("本程序用于把现有数据库中 dirs 和 files 的 id 、 parent_id 从 16 位截短到 8 位")
	println("您确定要执行这个操作吗？请输入 yes 或 no ：")

	CheckConfirm()
}
