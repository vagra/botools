package app

import (
	"fmt"
	"sync"
	"time"
)

func CheckSum() error {
	println("start: checksum of files")

	println()
	InitLog(CHECKSUM_LOG)

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
	CheckAllDBRootPathCorrect()

	println()
	GetNeedCheckSumDBs()

	CheckTaskHasDBs()
	InitMaps()

	println()
	MTCheckSum()

	println()
	println("checksum done!")
	return nil
}

func MTCheckSum() {
	println("每个 disk 启动多个 checker 和一个 writer 并行获取 sha1 和写入数据库")

	var main_wg sync.WaitGroup

	for name := range g_dbs {
		InitMap(name)
		GetCheckSumTasks(name)
	}

	for name := range g_dbs {
		if HasCheckSumTasks(name) {
			main_wg.Add(1)
			go CheckSumWorker(&main_wg, name)
		}
	}

	main_wg.Wait()
}

func CheckSumWorker(main_wg *sync.WaitGroup, disk_name string) {

	defer main_wg.Done()

	disk_path := g_disks[disk_name]

	fmt.Printf("%s worker: start scan %s\n", disk_name, disk_path)

	var wg sync.WaitGroup

	inChan := make(chan *File, MAX_CHAN)
	outChan := make(chan *File, MAX_CHAN)

	start := time.Now()

	wg.Add(1)
	go SumWriter(&wg, disk_name, outChan)

	for i := 0; i < g_threads; i++ {
		wg.Add(1)
		go SumChecker(&wg, disk_name, i, inChan, outChan)
	}

	for _, file := range g_map_files[disk_name] {
		inChan <- file
	}

	for i := 0; i < g_threads; i++ {
		inChan <- nil
	}

	wg.Wait()

	fmt.Printf("%s worker: stop. times: %v\n", disk_name, time.Since(start))
}

func SumChecker(wg *sync.WaitGroup, disk_name string, i int, ci <-chan *File, co chan<- *File) {
	defer wg.Done()

	for {
		select {

		case file := <-ci:

			if file == nil {
				fmt.Printf("%s: checker %d stop.\n", disk_name, i)

				co <- nil

				return
			}

			sha1, code := GetSHA1(file.path)
			file.sha1 = sha1
			file.status = code

			co <- file

		default:
			time.Sleep(time.Millisecond)
		}
	}
}

func SumWriter(wg *sync.WaitGroup, disk_name string, co <-chan *File) {
	defer wg.Done()

	var db *DB = g_dbs[disk_name]

	total := len(g_map_files[disk_name])

	var files []*File = []*File{}

	stops := 0
	count := 0
	for {
		select {

		case file := <-co:

			if file == nil {
				stops++
				if stops < g_threads {
					continue
				}

				fmt.Printf("%s: %d/%d\n", disk_name, count, total)

				db.BulkModFilesSha1(&files)
				files = nil

				fmt.Printf("%s: writer stop.\n", disk_name)

				return
			}

			files = append(files, file)

			count++
			if count%INSERT_COUNT == 0 {
				fmt.Printf("%s: %d/%d\n", disk_name, count, total)

				db.BulkModFilesSha1(&files)
				files = nil
			}

		default:
			time.Sleep(time.Millisecond)
		}
	}
}

func GetCheckSumTasks(disk_name string) {

	var db_path string = GetDBPath(disk_name)
	var db *DB = g_dbs[disk_name]

	fmt.Printf("%s: 从数据库 %s 中读取没有 sha1 的文件\n", disk_name, db_path)

	files := db.GetNoSHA1Files()

	g_map_files[disk_name] = files

	count := len(g_map_files[disk_name])

	if count == 0 {
		fmt.Printf("%s: 所有文件都已经有 sha1，不再重复获取。\n", disk_name)
		return
	}

	fmt.Printf("%s: 有 %d 个文件需要获取 sha1\n", disk_name, count)
}

func HasCheckSumTasks(disk_name string) bool {
	return len(g_map_files[disk_name]) > 0
}
