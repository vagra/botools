package app

import (
	"fmt"
	"sync"
)

func MoveLong() error {
	println("start: move longname dirs and files to special folder")

	println()
	InitLog(MOVE_LONG_LOG)

	println()
	ReadConfig()

	println()
	ReadErrors()

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
	ConfirmMoveLong()

	println()
	MTMoveLong()

	println()
	println("move longname dirs and files done!")
	return nil
}

func MTMoveLong() {
	println("每个 disk 启动一个线程，移动长名文件或文件夹，然后把 dirs 和 files 的 status 设为 3")

	var wg sync.WaitGroup

	for name := range g_dbs {
		wg.Add(1)
		go MoveLongWorker(&wg, name)
	}

	wg.Wait()
}

func MoveLongWorker(wg *sync.WaitGroup, disk_name string) {
	defer wg.Done()

	disk_path := g_disks[disk_name]

	fmt.Printf("%s worker: start scan %s\n", disk_name, disk_path)

	fmt.Printf("%s worker: stop.\n", disk_name)
}

func ConfirmMoveLong() {
	println("本程序用于同步物理目录到数据库")
	println("您确定要执行这个操作吗？请输入 yes 或 no ：")

	CheckConfirm()
}
