package app

import (
	"fmt"
	"sync"
)

func SyncVDB2Vir() error {
	println("start: sync vdb to vir tree")

	println()
	InitLog(VDB2VIR_LOG)

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
	ConfirmVDB2Vir()

	println()
	MTVDB2Vir()

	println()
	println("sync vdb to vir tree done!")

	return nil
}

func MTVDB2Vir() {
	println("每个 disk 启动一个线程，根据数据库中的 vdirs 和 vfiles 同步更新虚拟目录树")

	var wg sync.WaitGroup

	for name := range g_dbs {
		wg.Add(1)
		go VDB2VirWorker(&wg, name)
	}

	wg.Wait()
}

func VDB2VirWorker(wg *sync.WaitGroup, disk_name string) {
	defer wg.Done()

	disk_path := g_disks[disk_name]

	fmt.Printf("%s worker: start scan %s\n", disk_name, disk_path)

	fmt.Printf("%s worker: stop.\n", disk_name)
}

func ConfirmVDB2Vir() {
	println("本程序用于根据数据库中的 vdirs 和 vfiles 同步更新虚拟目录树")
	println("您确定要执行这个操作吗？请输入 yes 或 no ：")

	CheckConfirm()
}
