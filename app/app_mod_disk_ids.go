package app

import (
	"fmt"
	"sync"
)

func ModDiskIDs() error {
	println("start: change disk-id in dirs and files ids")

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
	CheckAllDBRootPathCorrect()

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
	println("change disk-id in db done!")
	return nil
}

func MTModDiskIDs() {
	println("每个 disk 启动一个线程，更新 dirs 和 files 的 id 和 parent_id 中的 disk-id")

	var wg sync.WaitGroup

	for name := range g_dbs {
		wg.Add(1)
		go ModDiskIDsWorker(&wg, name)
	}

	wg.Wait()
}

func ModDiskIDsWorker(wg *sync.WaitGroup, disk_name string) {
	defer wg.Done()

	var db *DB = g_dbs[disk_name]

	db_path := GetDBPath(disk_name)

	fmt.Printf("%s worker: start change disk-id in db %s\n", disk_name, db_path)

	dir_src := db.GetDirIDPrefix()
	dir_dst := GetDirPrefix(disk_name)

	file_src := db.GetFileIDPrefix()
	file_dst := GetFilePrefix(disk_name)

	fmt.Printf("%s -> %s\n", dir_src, dir_dst)
	fmt.Printf("%s -> %s\n", file_src, file_dst)

	db.ModDirsDiskID(dir_src, dir_dst)
	db.ModFilesDiskID(file_src, file_dst)

	fmt.Printf("%s worker: stop.\n", disk_name)
}

func ConfirmModDiskIDs() {
	println("本程序用于更改现有数据库中 dirs 和 files 的 id 、 parent_id 中的 disk-id")
	println("请确认已做好如下准备工作：")
	fmt.Printf("1. %s 中的 disks 已经更改为新的 disk-name\n", CONFIG_INI)
	fmt.Printf("2. %s 中的数据库文件也已经一一对应更改为新的 disk-name\n", DB_DIR)
	println("您确定要执行这个操作吗？请输入 yes 或 no ：")

	CheckConfirm()
}
