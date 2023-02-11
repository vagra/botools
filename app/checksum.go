package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

func CheckSum() error {
	println("start: checksum of files")

	if !AllDBExist() {
		println("数据库文件缺失，请重启本程序并选择 1 以初始化数据库")
		WaitExit(1)
	}

	if !ReadConfig() {
		WaitExit(1)
	}

	file, err := os.OpenFile(CHECKSUM_LOG, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	Check(err, "打开 "+CHECKSUM_LOG+" 时出错")
	defer file.Close()

	log.SetOutput(file)

	ReadSQL()
	GetDBs()

	if !HasData() {
		println("数据库中没有数据，请重启本程序并选择 2 以初始化数据库")
		WaitExit(1)
	}

	for name := range g_disks {
		GetDirs(name)
		if GetTasks(name) {
			GetSHA1(name)
		}
		println()
	}

	println("checksum done!")

	return nil
}

func GetDirs(disk_name string) {
	g_map_dirs = make(map[string]*Dir)
	var db_name string = GetDBName(disk_name)
	var db *sql.DB = g_dbs[db_name]

	fmt.Printf("从数据库 %s 中读取所有的目录\n", db_name)

	rows, err := g_dot.Query(db, SQL_GET_ALL_DIRS)
	Check(err, "执行 SQL "+SQL_GET_ALL_DIRS+" 时出错")
	defer rows.Close()

	for rows.Next() {
		var dir Dir

		err = rows.Scan(&dir.id, &dir.parent_id, &dir.name)
		Check(err, "执行 SQL "+SQL_GET_ALL_DIRS+" 后获取 dir 时出错")

		if dir.parent_id == "0" {
			dir.path = dir.name
		}

		g_map_dirs[dir.id] = &dir
	}

	CompDirsPath()
}

func GetTasks(disk_name string) bool {

	g_map_files = make(map[string]*File)

	var db_name string = GetDBName(disk_name)
	var db *sql.DB = g_dbs[db_name]

	fmt.Printf("从数据库 %s 中读取所有尚未获取 sha1 的文件\n", db_name)

	rows, err := g_dot.Query(db, SQL_GET_FILES_NO_SHA1)
	Check(err, "执行 SQL "+SQL_GET_FILES_NO_SHA1+" 时出错")
	defer rows.Close()

	for rows.Next() {
		var file File

		err = rows.Scan(&file.id, &file.parent_id, &file.name)
		Check(err, "执行 SQL "+SQL_GET_FILES_NO_SHA1+" 后获取 file 时出错")

		g_map_files[file.id] = &file
	}

	count := len(g_map_files)

	if count == 0 {
		println("所有文件都已经有了 sha1 值，不再重复获取。")
		return false
	}

	fmt.Printf("有 %d 个文件需要获取 sha1 值\n", count)

	CompFilesPath()

	return true
}

func GetSHA1(disk_name string) {

	fmt.Printf("启动多线程获取 %s 的文件 sha1 、单线程更新数据库...\n", disk_name)

	var wg sync.WaitGroup

	inChan := make(chan *File, 100)
	outChan := make(chan *File, 100)
	endChan := make(chan bool)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var db_name string = GetDBName(disk_name)
	var db *sql.DB = g_dbs[db_name]

	start := time.Now()

	go Writer(&wg, ctx, db, outChan, endChan)

	for i := 0; i < g_threads; i++ {
		go Checker(&wg, ctx, i, inChan, outChan)
	}

	for _, file := range g_map_files {
		inChan <- file
	}
	inChan <- nil

	end := <-endChan
	if end {
		println("endChan -> main: no more files.")
		println("main -> ctx: everyone stop!")
		cancel()
	}

	wg.Wait()

	fmt.Printf("执行时间: %v\n", time.Since(start))
}

func CompDirsPath() {
	println("合成目录路径")

	var count int64 = 0
	for {
		count = GetNoPathDirsCount()
		if count == 0 {
			break
		}

		DirsAddParentPath()
	}
}

func CompFilesPath() {
	println("合成文件路径")

	FilesAddParentPath()

}

func DirsAddParentPath() {

	for _, dir := range g_map_dirs {
		if dir.parent_id == "0" {
			continue
		}

		if len(dir.path) > 0 {
			continue
		}

		parent_path := g_map_dirs[dir.parent_id].path

		if len(parent_path) == 0 {
			continue
		}

		dir.path = parent_path + "/" + dir.name
	}
}

func FilesAddParentPath() {

	for _, file := range g_map_files {

		if len(file.path) > 0 {
			continue
		}

		parent_path := g_map_dirs[file.parent_id].path

		if len(parent_path) == 0 {
			continue
		}

		file.path = parent_path + "/" + file.name
	}
}

func GetNoPathDirsCount() int64 {

	var count int64 = 0
	for _, dir := range g_map_dirs {
		if len(dir.path) == 0 {
			count++
		}
	}

	return count
}
