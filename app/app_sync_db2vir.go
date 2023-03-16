package app

import (
	"fmt"
	"sync"
	"time"
)

func SyncDB2Vir() error {
	println("start: sync db to vir tree")

	println()
	InitLog(DB2VIR_LOG)

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
	CheckVirsRootExists()

	println()
	CheckErrorsRootExists()

	println()
	ConfirmDB2Vir()

	println()
	OnlyReadHasDataDBs2Mem()

	CheckTaskHasDBs()

	println()
	MTDB2Vir()

	println()
	println("sync db to vir tree done!")

	return nil
}

func MTDB2Vir() {
	println("每个 disk 启动一个线程，根据数据库中的 dirs 和 files 生成虚拟目录树")

	var wg sync.WaitGroup

	for name := range g_dbs {
		wg.Add(1)
		go DB2VirWorker(&wg, name)
	}

	wg.Wait()
}

func DB2VirWorker(wg *sync.WaitGroup, disk_name string) {
	defer wg.Done()

	disk_path := g_disks[disk_name]
	fmt.Printf("%s worker: start scan %s\n", disk_name, disk_path)

	start := time.Now()

	DB2Vir(disk_name)

	fmt.Printf("%s worker: stop. times: %v\n", disk_name, time.Since(start))
}

func ConfirmDB2Vir() {
	println("本程序用于根据数据库中的 dirs 和 files 在 virs-root 下生成虚拟目录树")
	println("请确认已做好如下准备工作：")
	fmt.Printf("1. 在 %s 中设置好 virs-root\n", CONFIG_INI)
	println("2. virs-root 目录若不为空，请手动清空")
	println("   因为程序将重新生成整个虚拟目录树，如果此前有多余的虚拟文件，不会被自动删除")
	fmt.Printf("3. 在 %s 中设置好 errors-root\n", CONFIG_INI)
	println("   异常文件的链接将会根据 errors-root 生成")
	println("您确定要执行这个操作吗？请输入 yes 或 no ：")

	CheckConfirm()
}
