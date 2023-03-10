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
	println("start: get tree")

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
	GetEmptyDBs()

	CheckTaskHasDBs()
	InitMaps()

	println()
	MTGetTree()

	println()
	println("get tree done!")
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

	disk_path := g_disks[disk_name]
	fmt.Printf("%s worker: start scan %s\n", disk_name, disk_path)

	start := time.Now()

	InitMap(disk_name)
	InitRootDir(disk_name, disk_path)

	ReadTree(disk_name)
	ReportDiskCounts(disk_name, disk_path)

	db_path := GetDBPath(disk_name)
	fmt.Printf("%s worker: write to db %s\n", disk_name, db_path)

	WriteDB(disk_name)
	ReportDBCounts(disk_name, db_path)

	fmt.Printf("%s worker: stop. times: %v\n", disk_name, time.Since(start))
}

func InitRootDir(disk_name string, disk_path string) {
	var dir Dir
	dir.id = GenDirUID(disk_name)
	dir.parent_id = "0"
	dir.name = disk_path
	dir.path = disk_path

	g_map_dirs[disk_name][dir.id] = &dir
}

func ReadTree(disk_name string) {
	root_id := GetDirUID(disk_name, 1)
	root_dir := g_map_dirs[disk_name][root_id]

	ReadDir(disk_name, root_dir, root_dir.name)
}

func WriteDB(disk_name string) {
	InsertDirs(disk_name)
	InsertFiles(disk_name)
}

func ReportDiskCounts(disk_name string, disk_path string) {
	fmt.Printf("%s %s\n%8d dirs, %8d files\n",
		disk_name, disk_path,
		len(g_map_dirs[disk_name]), len(g_map_files[disk_name]))
}

func ReportDBCounts(disk_name string, db_path string) {
	var db *DB = g_dbs[disk_name]

	fmt.Printf("%s %s\n%8d dirs, %8d files\n",
		disk_name, db_path,
		db.QueryDirsCount(), db.QueryFilesCount())
}

func ReadDir(disk_name string, dir *Dir, path string) {

	if !DirExists(path) {
		log.Printf("dir not exists: %s\n", path)
		return
	}

	items, _ := os.ReadDir(path)
	for _, item := range items {
		item_path := path + "/" + item.Name()
		item_path = strings.Replace(item_path, "//", "/", -1)

		if IsHidden(item_path) {
			continue
		}

		if item.IsDir() {

			var sub Dir
			sub.id = GenDirUID(disk_name)
			sub.parent_id = dir.id
			sub.name = item.Name()
			sub.path = item_path
			info, _ := item.Info()
			sub.mod_time = info.ModTime().Format(TIME_FORMAT)

			g_map_dirs[disk_name][sub.id] = &sub

			ReadDir(disk_name, &sub, item_path)

		} else {

			var file File
			file.id = GenFileUID(disk_name)
			file.parent_id = dir.id
			file.name = item.Name()
			file.path = item_path
			info, _ := item.Info()
			file.size = info.Size()
			file.mod_time = info.ModTime().Format(TIME_FORMAT)

			g_map_files[disk_name][file.id] = &file
		}
	}
}

func InsertDirs(disk_name string) {

	var db *DB = g_dbs[disk_name]

	var m int = 0
	var n int = 0

	db.BeginBulk()

	for _, dir := range g_map_dirs[disk_name] {

		db.AddDir(dir)

		m += 1
		n += 1

		if m >= len(g_map_dirs[disk_name]) {
			db.EndBulk()
			break
		}

		if n >= INSERT_COUNT {
			n = 0

			db.EndBulk()
			db.BeginBulk()
		}
	}
}

func InsertFiles(disk_name string) {
	var db *DB = g_dbs[disk_name]

	var m int = 0
	var n int = 0

	db.BeginBulk()

	for _, file := range g_map_files[disk_name] {

		db.AddFile(file)

		m += 1
		n += 1

		if m >= len(g_map_files[disk_name]) {
			db.EndBulk()
			break
		}

		if n >= INSERT_COUNT {
			n = 0

			db.EndBulk()
			db.BeginBulk()
		}
	}
}
