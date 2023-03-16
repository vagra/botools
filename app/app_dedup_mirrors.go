package app

import (
	"fmt"
	"sync"
	"time"
)

func DedupMirrors() error {
	println("start: de duplications in mirror")

	println()
	InitLog(DEDUP_MIRRORS_LOG)

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
	CheckAllDBRootPathInDisksRoot()

	println()
	ConfirmDedupMirrors()

	println()
	OnlyReadHasDataDBs2Mem()

	CheckTaskHasDBs()

	println()
	MTDedupMirrors()

	println()
	println("de duplications in mirror done!")

	return nil
}

func MTDedupMirrors() {
	println("每个 disk 启动一个线程，根据已查重的数据库，删除 mirrors-root 下的重复文件")

	var wg sync.WaitGroup

	for name := range g_dbs {
		wg.Add(1)
		go DedupMirrorsWorker(&wg, name)
	}

	wg.Wait()
}

func DedupMirrorsWorker(wg *sync.WaitGroup, disk_name string) {

	defer wg.Done()

	disk_path := g_disks[disk_name]
	fmt.Printf("%s worker: start dedup in %s\n", disk_name, disk_path)

	start := time.Now()

	DedupMirror(disk_name)

	fmt.Printf("%s worker: stop. times: %v\n", disk_name, time.Since(start))
}

func ConfirmDedupMirrors() {
	println("本程序用于根据已查重的数据库，在镜像目录下删除所有的重复文件")
	println("请确认已做好如下准备工作：")
	println("1. 将整个 disks-root 复制或建立硬链接到 mirrors-root")
	fmt.Printf("2. 在 %s 中设置好 disks-root 和 mirrors-root\n", CONFIG_INI)
	fmt.Printf("3. 检查数据库中的 dirs 和 files 根目录与 disks-root 一致\n")
	println("   程序通过把文件路径中的 disks-root 替换为 mirrors-root 来获得它在镜像目录下的路径\n")
	println("您确定要执行这个操作吗？请输入 yes 或 no ：")

	CheckConfirm()
}
