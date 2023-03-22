package app

import (
	"fmt"
	"sync"
)

func ResetDupIDs() error {
	println("start: clean files dup_id in db")

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
	ConfirmResetDupIDs()

	println()
	MTResetDupIDs()

	println()
	println("chean files dup_id in db done!")
	return nil
}

func MTResetDupIDs() {
	println("每个 disk 启动一个线程，重置 files 的 dup_id")

	var wg sync.WaitGroup

	for name := range g_dbs {
		wg.Add(1)
		go ResetDupIDsWorker(&wg, name)
	}

	wg.Wait()
}

func ResetDupIDsWorker(wg *sync.WaitGroup, disk_name string) {
	defer wg.Done()

	var db *DB = g_dbs[disk_name]

	db_path := GetDBPath(disk_name)

	fmt.Printf("%s worker: start clean files dup_id in db %s\n", disk_name, db_path)

	db.ResetFilesDupID()

	fmt.Printf("%s worker: stop.\n", disk_name)
}

func ConfirmResetDupIDs() {
	println("本程序用于清空数据库中的 files 的 dup_id")
	println("希望重新 dedup_dbs 和 dedup_mirrors 时需要这一步")
	println("您确定要执行这个操作吗？请输入 yes 或 no ：")

	CheckConfirm()
}
