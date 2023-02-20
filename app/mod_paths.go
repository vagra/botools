package app

import (
	"database/sql"
	"fmt"
	"sync"
)

func ModPaths() error {
	println("start: replace paths in db with new disk path")

	CheckConfig()

	println()
	GetDBs()
	ReadSQL()

	println()
	CheckAllDBExist()

	if !ConfirmModPaths() {
		WaitExit(1)
	}

	MTModPaths()

	println()
	println("replace paths in db done!")
	return nil
}

func MTModPaths() {
	println("每个 disk 启动一个线程，替换 dirs 和 files 的 path")

	var wg sync.WaitGroup

	for name, path := range g_disks {
		print(path)
		wg.Add(1)
		go ModPathsWorker(&wg, name, path)
	}

	wg.Wait()
}

func ModPathsWorker(wg *sync.WaitGroup, disk_name string, disk_path string) {
	defer wg.Done()

	var db *sql.DB = g_dbs[disk_name]

	fmt.Printf("%s worker: start replace paths %s\n", disk_name, GetDBPath(disk_name))

	root_dir := DBGetRootDir(db)

	DBUpdateRootDir(db, disk_path)

	DBReplaceDirPaths(db, root_dir.path, disk_path)
	DBReplaceFilePaths(db, root_dir.path, disk_path)

	fmt.Printf("%s worker: stop.\n", disk_name)
}

func ConfirmModPaths() bool {
	println()
	println("本程序用于把现有数据库中 dirs 和 files 的 path 替换为 disk 的新路径")
	println("您确定要执行这个操作吗？请输入 yes 或 no ：")

	return Confirm()
}
