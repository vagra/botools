package app

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
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

	InitMap(disk_name)
	InitRootDir(disk_name, disk_path)
	ReadTree(disk_name)
	WriteDB(disk_name)
	ReportCounts(disk_name, disk_path)

	fmt.Printf("%s worker: stop.\n", disk_name)
}

func InitRootDir(disk_name string, disk_path string) {
	var dir Dir
	dir.id = GenUID(disk_name, g_dirs_counter[disk_name])
	dir.parent_id = "0"
	dir.name = disk_path
	dir.path = disk_path

	g_map_dirs[disk_name][dir.id] = &dir
}

func ReadTree(disk_name string) {
	root_id := GetUID(disk_name, 1)
	root_dir := g_map_dirs[disk_name][root_id]

	ReadDir(disk_name, root_dir, root_dir.name)
}

func WriteDB(disk_name string) {
	InsertDirs(disk_name, INSERT_COUNT)
	InsertFiles(disk_name, INSERT_COUNT)
}

func ReportCounts(disk_name string, disk_path string) {
	var db *sql.DB = g_dbs[disk_name]

	fmt.Printf("%s %s\n", disk_name, disk_path)
	fmt.Printf("mem dirs: %d \t mem files: %d \n",
		len(g_map_dirs[disk_name]), len(g_map_files[disk_name]))
	fmt.Printf(" db dirs: %d \t  db files: %d \n",
		DBQueryDirsCount(db), DBQueryFilesCount(db))
}

func ReadDir(disk_name string, dir *Dir, path string) {

	if !DirExists(path) {
		log.Printf("dir not exists: %s\n", path)
		return
	}

	items, _ := ioutil.ReadDir(path)
	for _, item := range items {
		item_path := path + "/" + item.Name()
		item_path = strings.Replace(item_path, "//", "/", -1)

		if IsHidden(item_path) {
			continue
		}

		if item.IsDir() {

			var sub Dir
			sub.id = GenUID(disk_name, g_dirs_counter[disk_name])
			sub.parent_id = dir.id
			sub.name = item.Name()
			sub.path = item_path
			sub.mod_time = item.ModTime().Format(TIME_FORMAT)

			g_map_dirs[disk_name][sub.id] = &sub

			ReadDir(disk_name, &sub, item_path)

		} else {

			var file File
			file.id = GenUID(disk_name, g_files_counter[disk_name])
			file.parent_id = dir.id
			file.name = item.Name()
			file.path = item_path
			file.size = item.Size()
			file.mod_time = item.ModTime().Format(TIME_FORMAT)

			g_map_files[disk_name][file.id] = &file
		}
	}
}

func InsertDirs(disk_name string, count int) {

	var db *sql.DB = g_dbs[disk_name]

	var marks []string = []string{}
	var args []interface{} = []interface{}{}

	var m int = 0
	var n int = 0

	for _, dir := range g_map_dirs[disk_name] {
		m += 1
		n += 1

		dir.AddMarks(&marks)
		dir.AddArgs(&args)

		if n >= count || m >= len(g_map_dirs[disk_name]) {

			stmt := g_dot.QueryMap()[SQL_ADD_DIRS] + strings.Join(marks, ",\n")

			_, err := db.Exec(stmt, args...)
			Check(err, "在 dirs 表批量插入数据失败")

			marks = []string{}
			args = []interface{}{}
			n = 0
		}
	}
}

func InsertFiles(disk_name string, count int) {
	var db *sql.DB = g_dbs[disk_name]

	var marks []string = []string{}
	var args []interface{} = []interface{}{}

	var m int = 0
	var n int = 0

	for _, dir := range g_map_files[disk_name] {
		m += 1
		n += 1

		dir.AddMarks(&marks)
		dir.AddArgs(&args)

		if n >= count || m >= len(g_map_files[disk_name]) {

			stmt := g_dot.QueryMap()[SQL_ADD_FILES] + strings.Join(marks, ",\n")

			_, err := db.Exec(stmt, args...)
			Check(err, "在 files 表批量插入数据失败")

			marks = []string{}
			args = []interface{}{}
			n = 0
		}
	}
}