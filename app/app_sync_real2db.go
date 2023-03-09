package app

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
)

func SyncReal2DB() error {
	println("start: sync real tree to db")

	println()
	InitLog(REAL2DB_LOG)

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
	GetHasDataDBs()

	CheckTaskHasDBs()

	println()
	ConfirmReal2DB()

	println()
	MTReal2DB()

	println()
	println("sync real tree to db done!")

	return nil
}

func MTReal2DB() {
	println("每个 disk 启动一个线程，检查物理目录，更新数据库中 dirs 和 files 的 status")

	var wg sync.WaitGroup

	for name := range g_dbs {
		wg.Add(1)
		go Real2DBWorker(&wg, name)
	}

	wg.Wait()
}

func Real2DBWorker(wg *sync.WaitGroup, disk_name string) {

	defer wg.Done()

	disk_path := g_disks[disk_name]

	fmt.Printf("%s worker: start scan %s\n", disk_name, disk_path)

	src_path := GetDBPath(disk_name)
	dst_path := strings.Replace(src_path, ".db", ".new.db", -1)

	src_db := g_dbs[disk_name]

	mem_path := fmt.Sprintf("file:%s?mode=memory&cache=shared", disk_name)
	mem_db, err := sql.Open("sqlite3", mem_path)
	Check(err, "open memory db %s failed", mem_path)

	dst_db, err := sql.Open("sqlite3", dst_path)
	Check(err, "open new db %s failed", dst_path)

	err = BackupDB(mem_db, src_db)
	Check(err, "error when read db %s to memory", src_path)

	err = BackupDB(dst_db, mem_db)
	Check(err, "error when backup memory db to %s", dst_path)

	fmt.Printf("%s worker: stop.\n", disk_name)
}

func ConfirmReal2DB() {
	println("本程序用于同步物理目录到数据库中的 dirs 和 files 表")
	println("1. 如果目录或文件不存在，将其 status 设为 1，反之则设为 0，但不删除条目")
	println("2. 如果有新增目录或文件，将其插入 dirs 或 files 表")
	println("您确定要执行这个操作吗？请输入 yes 或 no ：")

	CheckConfirm()
}
