package app

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func GetTree() error {
	println("start: get tree")

	file, err := os.OpenFile(GET_TREE_LOG, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	Check(err, "打开 "+GET_TREE_LOG+" 时出错")
	defer file.Close()

	log.SetOutput(file)

	if !ReadConfig() {
		WaitExit(1)
	}

	CheckDBsDir()

	println()
	GetDBs()
	ReadSQL()

	if !AllDBExist() {
		println()
		println("检查到一些 disk 还没有数据库，请重启本程序并选择 1 以初始化数据库")
		println("或者修改 " + CONFIG_INI + " 用 # 注释掉不需要处理的 disk")
		WaitExit(1)
	}

	println()
	if HasData() {
		println("检查到一些数据库中存在数据，为避免重复生成数据，请重启本程序并选择 1 以初始化数据库")
		println("或者修改 " + CONFIG_INI + " 用 # 注释掉不需要处理的 disk")
		WaitExit(1)
	}

	for name, path := range g_disks {
		InitMaps(name, path)
		ReadTree(name)
		WriteDB(name)
		QueryCount(name)
		println()
	}

	println("get tree done!")
	return nil
}

func HasData() bool {

	for db_name, db := range g_dbs {

		if QueryDirsCount(db) > 0 {
			println("数据库 " + db_name + " 的 dirs 表中存在数据")
			return true
		}

		if QueryFilesCount(db) > 0 {
			println("数据库 " + db_name + " 的 files 表中存在数据")
			return true
		}
	}

	return false
}

func InitMaps(disk_name string, disk_path string) {
	g_map_dirs = make(map[string]Dir)
	g_map_files = make(map[string]File)

	g_dirs_counter = 0
	g_files_counter = 0

	var dir Dir
	dir.id = GenUID(disk_name, &g_dirs_counter)
	dir.name = disk_path
	dir.parent_id = "0"

	g_map_dirs[dir.id] = dir
}

func ReadTree(disk_name string) {
	root_id := GetUID(disk_name, 1)
	root_dir := g_map_dirs[root_id]

	println("遍历 " + disk_name + ": " + root_dir.name)
	ReadDir(disk_name, root_dir, root_dir.name)
}

func WriteDB(disk_name string) {

	var db *sql.DB = g_dbs[GetDBName(disk_name)]

	InsertDirs(db, INSERT_COUNT)
	InsertFiles(db, INSERT_COUNT)

}

func QueryCount(disk_name string) {
	var db *sql.DB = g_dbs[GetDBName(disk_name)]

	fmt.Printf("mem dirs: %d \t mem files: %d \n", len(g_map_dirs), len(g_map_files))
	fmt.Printf(" db dirs: %d \t  db files: %d \n", QueryDirsCount(db), QueryFilesCount(db))
}

func ReadDir(disk_name string, dir Dir, path string) {

	if !DirExist(path) {
		log.Println("dir not exist: " + path)
		return
	}

	items, _ := ioutil.ReadDir(path)
	for _, item := range items {
		item_path := path + "/" + item.Name()

		if IsHidden(item_path) {
			continue
		}

		if item.IsDir() {

			var sub Dir
			sub.id = GenUID(disk_name, &g_dirs_counter)
			sub.name = item.Name()
			sub.parent_id = dir.id
			sub.mod_time = item.ModTime().Format(TIME_FORMAT)

			g_map_dirs[sub.id] = sub

			ReadDir(disk_name, sub, item_path)

		} else {

			var file File
			file.id = GenUID(disk_name, &g_files_counter)
			file.name = item.Name()
			file.parent_id = dir.id
			file.size = item.Size()
			file.mod_time = item.ModTime().Format(TIME_FORMAT)

			g_map_files[file.id] = file
		}
	}
}

func InsertDirs(db *sql.DB, count int) {

	var marks []string = []string{}
	var args []interface{} = []interface{}{}

	var m int = 0
	var n int = 0

	for _, dir := range g_map_dirs {
		m += 1
		n += 1

		dir.AddMarks(&marks)
		dir.AddArgs(&args)

		if n >= count || m >= len(g_map_dirs) {

			stmt := SQL_ADD_DIRS + strings.Join(marks, ",\n")

			_, err := db.Exec(stmt, args...)
			Check(err, "在 dirs 表中批量插入数据失败")

			marks = []string{}
			args = []interface{}{}
			n = 0
		}
	}
}

func InsertFiles(db *sql.DB, count int) {

	var marks []string = []string{}
	var args []interface{} = []interface{}{}

	var m int = 0
	var n int = 0

	for _, dir := range g_map_files {
		m += 1
		n += 1

		dir.AddMarks(&marks)
		dir.AddArgs(&args)

		if n >= count || m >= len(g_map_files) {

			stmt := SQL_ADD_FILES + strings.Join(marks, ",\n")

			_, err := db.Exec(stmt, args...)
			Check(err, "在 files 表中批量插入数据失败")

			marks = []string{}
			args = []interface{}{}
			n = 0
		}
	}
}

func QueryDirsCount(db *sql.DB) int64 {

	var count int64 = 0

	row, err := g_dot.QueryRow(db, SQL_COUNT_DIRS)
	Check(err, "执行 SQL "+SQL_COUNT_DIRS+" 时出错")

	err = row.Scan(&count)
	Check(err, "执行 SQL "+SQL_COUNT_DIRS+" 后获取 count 时出错")

	return count
}

func QueryFilesCount(db *sql.DB) int64 {

	var count int64 = 0

	row, err := g_dot.QueryRow(db, SQL_COUNT_FILES)
	Check(err, "执行 SQL "+SQL_COUNT_FILES+" 时出错")

	err = row.Scan(&count)
	Check(err, "执行 SQL "+SQL_COUNT_FILES+" 后获取 count 时出错")

	return count
}
