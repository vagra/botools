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

	file, err := os.OpenFile(CHECKSUM_LOG, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	Check(err, "打开 "+CHECKSUM_LOG+" 时出错")
	defer file.Close()

	log.SetOutput(file)

	CheckConfig()

	println()
	ReadSQL()
	GetDBs()

	println()
	CheckAllDBExist()
	CheckAnyDBHasData()

	InitMaps()

	MTCheckSum()

	println()
	println("checksum done!")
	return nil
}

func MTCheckSum() {
	println("每个 disk 启动多个 checker 和一个 writer 并行获取 sha1 和写入数据库")

	var main_wg sync.WaitGroup

	for name := range g_disks {
		InitMap(name)
		GetTasks(name)
	}

	for name, path := range g_disks {
		if HasTasks(name) {
			main_wg.Add(1)
			go CheckSumWorker(&main_wg, name, path)
		}
	}

	main_wg.Wait()
}

func CheckSumWorker(main_wg *sync.WaitGroup, disk_name string, disk_path string) {

	defer main_wg.Done()

	fmt.Printf("%s worker: start scan %s\n", disk_name, disk_path)

	var wg sync.WaitGroup

	inChan := make(chan *File, MAX_CHAN)
	outChan := make(chan *File, MAX_CHAN)
	endChan := make(chan bool)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	start := time.Now()

	wg.Add(1)
	go Writer(&wg, ctx, disk_name, outChan, endChan)

	for i := 0; i < g_threads; i++ {
		wg.Add(1)
		go Checker(&wg, ctx, disk_name, i, inChan, outChan)
	}

	for _, file := range g_map_files[disk_name] {
		inChan <- file
	}
	inChan <- nil

	end := <-endChan
	if end {
		fmt.Printf("%s: endChan -> main: no more files.\n", disk_name)
		fmt.Printf("%s: main -> ctx: everyone stop!\n", disk_name)
		cancel()
	}

	wg.Wait()

	fmt.Printf("%s worker: stop. times: %v\n", disk_name, time.Since(start))
}

func GetDirs(disk_name string) {

	var db_path string = GetDBPath(disk_name)
	var db *sql.DB = g_dbs[disk_name]

	fmt.Printf("从数据库 %s 中读取所有的目录\n", db_path)

	g_map_dirs[disk_name] = DBGetAllDirs(db)
}

func GetTasks(disk_name string) {

	var db_path string = GetDBPath(disk_name)
	var db *sql.DB = g_dbs[disk_name]

	fmt.Printf("%s: 从数据库 %s 中读取所有尚未获取 sha1 的文件\n", disk_name, db_path)

	files := DBGetFilesNoSHA1(db)

	g_map_files[disk_name] = files

	count := len(g_map_files[disk_name])

	if count == 0 {
		fmt.Printf("%s: 所有文件都已经有了 sha1 值，不再重复获取。\n", disk_name)
		return
	}

	fmt.Printf("%s: 有 %d 个文件需要获取 sha1 值\n", disk_name, count)
}

func HasTasks(disk_name string) bool {
	return len(g_map_files[disk_name]) > 0
}

func CompDirsPath(disk_name string) {
	println("合成目录路径")

	var count int64 = 0
	for {
		count = GetNoPathDirsCount(disk_name)
		if count == 0 {
			break
		}

		DirsAddParentPath(disk_name)
	}
}

func CompFilesPath(disk_name string) {
	println("合成文件路径")

	FilesAddParentPath(disk_name)
}

func DirsAddParentPath(disk_name string) {

	for _, dir := range g_map_dirs[disk_name] {
		if dir.parent_id == "0" {
			continue
		}

		if len(dir.path) > 0 {
			continue
		}

		parent_path := g_map_dirs[disk_name][dir.parent_id].path

		if len(parent_path) == 0 {
			continue
		}

		dir.path = parent_path + "/" + dir.name
	}
}

func FilesAddParentPath(disk_name string) {

	for _, file := range g_map_files[disk_name] {

		if len(file.path) > 0 {
			continue
		}

		parent_path := g_map_dirs[disk_name][file.parent_id].path

		if len(parent_path) == 0 {
			continue
		}

		file.path = parent_path + "/" + file.name
	}
}

func GetNoPathDirsCount(disk_name string) int64 {

	var count int64 = 0
	for _, dir := range g_map_dirs[disk_name] {
		if len(dir.path) == 0 {
			count++
		}
	}

	return count
}
