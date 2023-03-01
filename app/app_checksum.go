package app

import (
	"context"
	"database/sql"
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

func Checker(wg *sync.WaitGroup, ctx context.Context, disk_name string, i int, ci <-chan *File, co chan<- *File) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s: checker %d stop.\n", disk_name, i)
			return

		case file := <-ci:

			if file == nil {
				fmt.Printf("%s: inChan -> checker %d -> outChan: no more files.\n", disk_name, i)

				co <- nil

				continue
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

func Writer(wg *sync.WaitGroup, ctx context.Context, disk_name string, co <-chan *File, ce chan<- bool) {
	defer wg.Done()

	var db *sql.DB = g_dbs[disk_name]

	total := len(g_map_files[disk_name])
	divisor := int(total / 20)

	count := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s: writer stop.\n", disk_name)
			return

		case file := <-co:

			if file == nil {
				fmt.Printf("%s: outChan -> writer -> endChan: no more files.\n", disk_name)
				ce <- true

				continue
			}

			DBUpdateFile(db, file)

			count++

			if count%divisor == 0 {
				fmt.Printf("%s: %d%%\n", disk_name, count*100/total+1)
			}

		default:
			time.Sleep(time.Millisecond)
		}
	}
}

func GetCheckSumTasks(disk_name string) {

	var db_path string = GetDBPath(disk_name)
	var db *sql.DB = g_dbs[disk_name]

	fmt.Printf("%s: 从数据库 %s 中读取没有 sha1 的文件\n", disk_name, db_path)

	files := DBGetNoSHA1Files(db)

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
