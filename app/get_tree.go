package app

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func GetTree() error {
	println("start: get tree")

	if !AllDBExist() {
		println("数据库文件缺失，请重启本程序并选择 1 以初始化数据库")
		WaitExit(1)
	}

	if !ReadConfig() {
		WaitExit(1)
	}

	file, err := os.OpenFile(GET_TREE_LOG, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	Check(err, "打开 "+GET_TREE_LOG+" 时出错")
	defer file.Close()

	log.SetOutput(file)

	ReadSQL()
	GetDBs()

	if HasData() {
		println("数据库中存在数据，为避免重复生成数据，请重启本程序并选择 1 以初始化数据库")
		WaitExit(1)
	}

	ReadDisks()
	ReadDirs()

	println("get tree done!")
	return nil
}

func HasData() bool {

	var row *sql.Row
	var err error
	var count int64

	for i := 0; i < DB_COUNT; i++ {
		println("检查数据库 " + g_db_names[i] + " 是否为空")

		row, err = g_dot.QueryRow(g_dbs[i], g_count_sqls[i])
		Check(err, "执行 SQL "+g_count_sqls[i]+" 时出错")

		err = row.Scan(&count)
		Check(err, "执行 SQL "+g_count_sqls[i]+" 后获取数据时出错")

		if count > 0 {
			fmt.Printf("数据库 %s 中存在 %d 条数据\n", g_db_names[i], count)
			return true
		}
	}

	return false
}

func ReadDisks() {
	for name, path := range g_disks {
		ReadDisk(name, path)
	}
}

func ReadDirs() {
	rows, err := g_dot.Query(g_disks_db, SQL_GET_ALL_DISKS)
	Check(err, "执行 SQL "+SQL_GET_ALL_DISKS+" 时出错")

	var disk Disk
	var dir Dir

	for rows.Next() {
		err := rows.Scan(&disk.id, &disk.name, &disk.path, &disk.dir_id)
		Check(err, "执行 SQL "+SQL_GET_ALL_DISKS+" 后获取数据时出错")

		dir.id = disk.dir_id
		dir.name = disk.path
		dir.parent_id = 0

		println("遍历 " + disk.name + ": " + disk.path)
		ReadDir(dir, dir.name)
	}

}

func ReadDisk(name string, path string) {

	log.Println(path)

	var disk Disk
	disk.name = name
	disk.path = path

	var dir Dir
	dir.name = path
	dir.parent_id = 0

	disk.dir_id = InsertDir(dir)

	InsertDisk(disk)
}

func ReadDir(dir Dir, path string) {

	if !DirExist(path) {
		log.Println("dir not exist: " + path)
		return
	}

	items, _ := ioutil.ReadDir(path)
	for _, item := range items {
		item_path := path + "/" + item.Name()

		if item.IsDir() {

			var sub Dir
			sub.name = item.Name()
			sub.parent_id = dir.id

			sub.id = InsertDir(sub)

			var meta DirMeta
			meta.dir_id = sub.id
			meta.size = item.Size()
			meta.mod_time = item.ModTime().Format(TIME_FORMAT)

			InsertDirMeta(meta)

			ReadDir(sub, item_path)

		} else {

			var file File
			file.name = item.Name()
			file.parent_id = dir.id

			file.id = InsertFile(file)

			var meta FileMeta
			meta.file_id = file.id
			meta.size = item.Size()
			meta.mod_time = item.ModTime().Format(TIME_FORMAT)

			InsertFileMeta(meta)
		}
	}
}

func ReadFile(file File, path string) {

	if !FileExist(path) {
		log.Println("file not exist: " + path)
		return
	}
}

func InsertDisk(disk Disk) int64 {
	res, err := g_dot.Exec(g_disks_db, SQL_ADD_DISK, disk.name, disk.path, disk.dir_id)
	Check(err, "执行 SQL "+SQL_ADD_DISK+" 时出错")

	id, err := res.LastInsertId()
	Check(err, "执行 SQL "+SQL_ADD_DISK+" 后获取 id 时出错")

	return id
}

func InsertDir(dir Dir) int64 {
	res, err := g_dot.Exec(g_dirs_db, SQL_ADD_DIR, dir.name, dir.parent_id)
	Check(err, "执行 SQL "+SQL_ADD_DIR+" 时出错")

	id, err := res.LastInsertId()
	Check(err, "执行 SQL "+SQL_ADD_DIR+" 后获取 id 时出错")

	return id
}

func InsertFile(file File) int64 {
	res, err := g_dot.Exec(g_files_db, SQL_ADD_FILE, file.name, file.parent_id)
	Check(err, "执行 SQL "+SQL_ADD_FILE+" 时出错")

	id, err := res.LastInsertId()
	Check(err, "执行 SQL "+SQL_ADD_FILE+" 后获取 id 时出错")

	return id
}

func InsertDirMeta(meta DirMeta) {
	_, err := g_dot.Exec(g_dir_metas_db, SQL_ADD_DIR_META, meta.dir_id, meta.size, meta.mod_time)
	Check(err, "执行 SQL "+SQL_ADD_DIR_META+" 时出错")
}

func InsertFileMeta(meta FileMeta) {
	_, err := g_dot.Exec(g_file_metas_db, SQL_ADD_FILE_META, meta.file_id, meta.size, meta.mod_time)
	Check(err, "执行 SQL "+SQL_ADD_FILE_META+" 时出错")
}
