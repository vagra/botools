package app

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

func GetTree() error {
	println("start: get tree")

	file, err := os.OpenFile(GET_TREE_LOG, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	Check(err, "打开 "+GET_TREE_LOG+" 时出错")
	defer file.Close()

	log.SetOutput(file)

	CheckConfig()

	println()
	GetDBs()
	ReadSQL()

	println()
	CheckAllDBExist()
	CheckAllDBNoData()

	InitMaps()

	MTGetTree()

	println()
	println("get tree done!")
	return nil
}

func MTGetTree() {
	println("每个 disk 启动一个线程，先获取目录树，然后批量写入数据库")

	var wg sync.WaitGroup

	for name, path := range g_disks {
		wg.Add(1)
		go GetTreeWorker(&wg, name, path)
	}

	wg.Wait()
}

func GetTreeWorker(wg *sync.WaitGroup, disk_name string, disk_path string) {
	defer wg.Done()

	fmt.Printf("%s worker: start scan %s\n", disk_name, disk_path)

	InitMap(disk_name)
	InitRootDir(disk_name, disk_path)
	ReadTree(disk_name)
	WriteDB(disk_name)
	ReportCount(disk_name, disk_path)

	fmt.Printf("%s worker: stop.\n", disk_name)
}

func CheckAllDBExist() {
	if !AllDBExist() {
		println()
		println("检查到一些 disk 还没有数据库，请重启本程序并选择 1 以初始化数据库")
		fmt.Printf("或者修改 %s 用 # 注释掉不需要处理的 disk\n", CONFIG_INI)
		WaitExit(1)
	}
}

func CheckAllDBNoData() {
	if HasData() {
		println("检查到一些数据库中存在数据，为避免重复生成数据，请重启本程序并选择 1 以初始化数据库")
		fmt.Printf("或者修改 %s 用 # 注释掉不需要处理的 disk\n", CONFIG_INI)
		WaitExit(1)
	}
}

func CheckAnyDBHasData() {
	if !HasData() {
		println("数据库中没有数据，请重启本程序并选择 2 以初始化数据库")
		WaitExit(1)
	}
}

func HasData() bool {

	for db_name, db := range g_dbs {

		if DBQueryDirsCount(db) > 0 {
			fmt.Printf("数据库 %s 的 dirs 表中存在数据\n", db_name)
			return true
		}

		if DBQueryFilesCount(db) > 0 {
			fmt.Printf("数据库 %s 的 files 表中存在数据\n", db_name)
			return true
		}
	}

	return false
}

func InitMaps() {
	g_map_dirs = make(map[string]map[string]*Dir)
	g_map_files = make(map[string]map[string]*File)

	g_dirs_counter = make(map[string]*int64)
	g_files_counter = make(map[string]*int64)
}

func InitMap(disk_name string) {
	g_map_dirs[disk_name] = make(map[string]*Dir)
	g_map_files[disk_name] = make(map[string]*File)

	var dirs_counter int64 = 0
	var files_counter int64 = 0

	g_dirs_counter[disk_name] = &dirs_counter
	g_files_counter[disk_name] = &files_counter
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

func ReportCount(disk_name string, disk_path string) {
	var db *sql.DB = g_dbs[disk_name]

	fmt.Printf("%s %s\n", disk_name, disk_path)
	fmt.Printf("mem dirs: %d \t mem files: %d \n",
		len(g_map_dirs[disk_name]), len(g_map_files[disk_name]))
	fmt.Printf(" db dirs: %d \t  db files: %d \n",
		DBQueryDirsCount(db), DBQueryFilesCount(db))
}

func ReadDir(disk_name string, dir *Dir, path string) {

	if !DirExist(path) {
		log.Printf("dir not exist: %s\n", path)
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
			Check(err, "在 dirs 表中批量插入数据失败")

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
			Check(err, "在 files 表中批量插入数据失败")

			marks = []string{}
			args = []interface{}{}
			n = 0
		}
	}
}
