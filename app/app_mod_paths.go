package app

import (
	"database/sql"
	"fmt"
	"sync"
)

func ModPaths() error {
	println("start: replace paths in db with new disk path")

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
	ConfirmModPaths()

	println()
	MTModPaths()

	println()
	println("replace paths in db done!")
	return nil
}

func MTModPaths() {
	println("每个 disk 启动一个线程，替换 dirs 和 files 的 path")

	var wg sync.WaitGroup

	for name := range g_dbs {
		wg.Add(1)
		go ModPathsWorker(&wg, name)
	}

	wg.Wait()
}

func ModPathsWorker(wg *sync.WaitGroup, disk_name string) {
	defer wg.Done()

	var db *sql.DB = g_dbs[disk_name]

	db_path := GetDBPath(disk_name)

	fmt.Printf("%s worker: start replace paths in db %s\n", disk_name, db_path)

	root_dir := DBGetRootDir(db)

	old_root := root_dir.path
	new_root := g_disks[disk_name]

	if new_root == old_root {
		fmt.Printf("%s worker: same path %s\n", disk_name, old_root)
		fmt.Printf("%s worker: do nothing.\n", disk_name)
		fmt.Printf("%s worker: stop.\n", disk_name)
		return
	}

	fmt.Printf("%s worker: old path %s\n", disk_name, old_root)
	fmt.Printf("%s worker: new path %s\n", disk_name, new_root)

	DBUpdateRootDir(db, new_root)

	DBReplaceDirPaths(db, old_root, new_root)
	DBReplaceFilePaths(db, old_root, new_root)

	fmt.Printf("%s worker: stop.\n", disk_name)
}

func ConfirmModPaths() {
	println("本程序用于把现有数据库中 dirs 和 files 的 path 替换为 disks 的新路径")
	println("您确定要执行这个操作吗？请输入 yes 或 no ：")

	CheckConfirm()
}
