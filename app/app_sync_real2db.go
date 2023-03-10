package app

import (
	"fmt"
	"sync"
)

func SyncReal2DB() error {
	println("start: sync real tree to db")

	println()
	InitLog(REAL2DB_LOG)

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
	CheckAllDBRootPathCorrect()

	println()
	LoadHasDataDBs2Mem()

	CheckTaskHasDBs()

	println()
	ConfirmReal2DB()

	println()
	MTReal2DB()

	println()
	BakeMemDBs()

	println()
	println("sync real tree to db done!")

	return nil
}

func MTReal2DB() {
	println("每个 disk 启动一个线程，检查物理目录，更新数据库中 dirs 和 files 的 status")

	var wg sync.WaitGroup

	for name := range g_dbs {
		wg.Add(1)
		go Real2DBWorker(&wg, name)
	}

	wg.Wait()
}

func Real2DBWorker(wg *sync.WaitGroup, disk_name string) {

	defer wg.Done()

	disk_path := g_disks[disk_name]

	fmt.Printf("%s worker: start scan %s\n", disk_name, disk_path)

	db := g_dbs[disk_name]

	fmt.Printf("%s files: %d\n", disk_name, db.QueryFilesCount())

	fmt.Printf("%s worker: stop.\n", disk_name)
}

func ConfirmReal2DB() {
	println("本程序用于同步物理目录到数据库中的 dirs 和 files 表")
	println("1. 如果目录或文件不存在，将其 status 设为 1，反之则设为 0，但不删除条目")
	println("2. 如果有新增目录或文件，将其插入 dirs 或 files 表")
	println("您确定要执行这个操作吗？请输入 yes 或 no ：")

	CheckConfirm()
}
