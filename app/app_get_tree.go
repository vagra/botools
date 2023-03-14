package app

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

func GetTree() error {
	println("start: get tree from real disks")

	println()
	InitLog(GET_TREE_LOG)

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
	ConfirmGetTree()

	println()
	LoadEmptyDBs2Mem()

	CheckTaskHasDBs()

	InitDBCounters()

	println()
	MTGetTree()

	println()
	println("get tree from real disks done!")

	return nil
}

func MTGetTree() {
	println("每个 disk 启动一个线程，先获取目录树，然后批量写入数据库")

	var wg sync.WaitGroup

	for name := range g_dbs {
		wg.Add(1)
		go GetTreeWorker(&wg, name)
	}

	wg.Wait()
}

func GetTreeWorker(wg *sync.WaitGroup, disk_name string) {
	defer wg.Done()

	start := time.Now()

	disk_path := g_disks[disk_name]

	fmt.Printf("%s worker: start scan %s\n", disk_name, disk_path)

	ReadTree(disk_name)

	BakeMemDB(disk_name)

	fmt.Printf("%s worker: stop. times: %v\n", disk_name, time.Since(start))
}

func ConfirmGetTree() {
	println("本程序用于遍历物理目录，在数据库中建立最初的 dirs 和 files 条目")
	fmt.Printf("1. 执行本程序前，需要确保 %s 中所有 disks 对应的数据库都已初始化\n", CONFIG_INI)
	println("2. 如果一个数据库中已经存在 dirs 或 files，会跳过这个数据库")
	println("您确定要执行这个操作吗？请输入 yes 或 no ：")

	CheckConfirm()
}

func ReadTree(disk_name string) {

	disk_path := g_disks[disk_name]

	var root_dir Dir
	root_dir.parent_id = "0"
	root_dir.name = disk_path
	root_dir.path = disk_path

	db := g_dbs[disk_name]

	ReadDir(db, disk_name, &root_dir)

	fmt.Printf("%s: %8d dirs, %8d files\n",
		disk_name, db.QueryDirsCount(), db.QueryFilesCount())
}

func ReadDir(db *DB, disk_name string, dir *Dir) {

	if !DirExists(dir.path) {
		log.Printf("real dir not exists: %s\n", dir.path)
		return
	}

	dir.id = GenDirUID(disk_name)
	db.AddDir(dir)

	items, _ := os.ReadDir(dir.path)
	for _, item := range items {
		item_path := dir.path + "/" + item.Name()
		item_path = strings.Replace(item_path, "//", "/", -1)

		if IsHidden(item_path) {
			continue
		}

		if item.IsDir() {

			var sub Dir
			sub.parent_id = dir.id
			sub.name = item.Name()
			sub.path = item_path
			info, _ := item.Info()
			sub.mod_time = info.ModTime().Format(TIME_FORMAT)

			ReadDir(db, disk_name, &sub)

		} else {

			var file File

			file.parent_id = dir.id
			file.name = item.Name()
			file.path = item_path
			info, _ := item.Info()
			file.size = info.Size()
			file.mod_time = info.ModTime().Format(TIME_FORMAT)

			ReadFile(db, disk_name, &file)

		}
	}
}

func ReadFile(db *DB, disk_name string, file *File) {

	if !FileExists(file.path) {
		log.Printf("real file not exists: %s\n", file.path)
		return
	}

	file.id = GenFileUID(disk_name)
	db.AddFile(file)
}
