package app

import (
	"fmt"
	"log"
	"os"
	"sync"
)

func SyncReal2DB() error {
	println("start: sync real tree to db")

	file, err := os.OpenFile(REAL2DB_LOG, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	Check(err, "打开 %s 时出错", REAL2DB_LOG)
	defer file.Close()

	log.SetOutput(file)

	ReadConfig()

	println()
	ReadDotSQL()

	println()
	GetAllDBs()

	InitMaps()

	println()
	MTGetTree()

	println()
	println("sync real tree to db done!")

	return nil
}

func MTReal2DB() {
	println("每个 disk 启动一个线程，先获取目录树，然后批量写入数据库")

	var wg sync.WaitGroup

	for name := range g_dbs {
		path := g_disks[name]

		wg.Add(1)
		go GetTreeWorker(&wg, name, path)
	}

	wg.Wait()
}

func Real2DBWorker(wg *sync.WaitGroup, disk_name string, disk_path string) {
	defer wg.Done()

	fmt.Printf("%s worker: start scan %s\n", disk_name, disk_path)

	InitMap(disk_name)
	InitRootDir(disk_name, disk_path)
	ReadTree(disk_name)
	WriteDB(disk_name)
	ReportCounts(disk_name, disk_path)

	fmt.Printf("%s worker: stop.\n", disk_name)
}
