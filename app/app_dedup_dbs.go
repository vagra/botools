package app

import (
	"fmt"
	"time"
)

func DedupDBs() error {
	println("start: de duplications in db")

	println()
	InitLog(DEDUP_DBS_LOG)

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
	ConfirmDedupDBs()

	println()
	LoadHasDataDBs2Mem()

	CheckTaskHasDBs()

	ReadUniqueMap()

	println()
	STDedupDBs()

	println()
	println("de duplications in db done!")

	return nil
}

func STDedupDBs() {
	println("使用单线程，对数据库中的 files 查重")

	for name := range g_dbs {
		DedupDBsWorker(name)
		println()
	}
}

func DedupDBsWorker(disk_name string) {

	db_path := GetDBPath(disk_name)
	fmt.Printf("%s worker: start scan %s\n", disk_name, db_path)

	start := time.Now()

	DedupDB(disk_name)

	BakeMemDB(disk_name)

	fmt.Printf("%s worker: stop. times: %v\n", disk_name, time.Since(start))
}

func ConfirmDedupDBs() {
	println("本程序用于对数据库中的 files 查重")
	println("1. 如果存在重复，则只保留一个主文件，将重复文件的 dup_id 设为主文件的 id，但不会删除条目")
	println("2. 这是基于所有 disks 的跨盘查重")
	println("3. 不检查或删除物理文件")
	println("4. 已经标记为重复文件的，下次查重不会再处理它")
	println("您确定要执行这个操作吗？请输入 yes 或 no ：")

	CheckConfirm()
}
