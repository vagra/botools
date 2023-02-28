package app

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

func CheckDBsDirExists() {
	fmt.Printf("检查数据库文件夹 %s 是否存在\n", DB_DIR)

	if !DirExists(DB_DIR) {
		fmt.Printf("数据库文件夹 %s 不存在，新建...\n", DB_DIR)

		err := os.Mkdir(DB_DIR, os.ModePerm)
		Check(err, "创建数据库文件夹 %s 时出错", DB_DIR)
	}
}

func GetDBPath(disk_name string) string {
	return fmt.Sprintf("%s/%s%s", DB_DIR, disk_name, DB_EXT)
}

func CheckTaskHasDBs() {
	if len(g_dbs) > 0 {
		return
	}

	println("没有需要处理的数据库")
	WaitExit(0)
}

func GetAllDBs() {
	g_dbs = make(map[string]*sql.DB)

	for disk_name := range g_disks {

		db_path := GetDBPath(disk_name)

		CheckDBExists(disk_name)

		fmt.Printf("打开数据库 %s\n", db_path)

		db := DBOpen(db_path)

		g_dbs[disk_name] = db
	}
}

func GetNotExistsDBs() {
	println("检查还不存在的数据库")

	g_dbs = make(map[string]*sql.DB)

	for disk_name := range g_disks {

		db_path := GetDBPath(disk_name)

		if DBExists(disk_name) {
			fmt.Printf("数据库 %s 已存在，跳过\n", db_path)
			continue
		}

		fmt.Printf("数据库 %s 不存在，新建...\n", db_path)

		db, err := sql.Open("sqlite3", db_path)
		Check(err, "打开数据库 %s 失败", db_path)

		g_dbs[disk_name] = db
	}
}

func GetEmptyDBs() {
	println("获取所有 dirs 和 files 表为空的数据库")

	g_dbs = make(map[string]*sql.DB)

	for disk_name := range g_disks {

		db_path := GetDBPath(disk_name)

		CheckDBExists(disk_name)

		db := DBOpen(db_path)

		CheckDBInited(db, db_path)

		if _, yes := DBNoData(db); !yes {
			continue
		}
		println(db_path)

		g_dbs[disk_name] = db
	}
}

func GetHasDataDBs() {
	println("获取所有 dirs 和 files 表都有数据的数据库")

	g_dbs = make(map[string]*sql.DB)

	for disk_name := range g_disks {

		db_path := GetDBPath(disk_name)

		CheckDBExists(disk_name)

		db := DBOpen(db_path)

		CheckDBInited(db, db_path)

		CheckDBHasData(db, db_path)

		println(db_path)

		g_dbs[disk_name] = db
	}
}

func GetNeedCheckSumDBs() {
	println("获取还有 files 需要获取 sha1 的数据库")

	g_dbs = make(map[string]*sql.DB)

	for disk_name := range g_disks {

		db_path := GetDBPath(disk_name)

		CheckDBExists(disk_name)

		db := DBOpen(db_path)

		CheckDBInited(db, db_path)

		CheckDBHasData(db, db_path)

		if _, yes := DBNeedCheckSum(db); !yes {
			continue
		}

		println(db_path)

		g_dbs[disk_name] = db
	}
}

func CheckDBExists(db_name string) {
	if !DBExists(db_name) {
		db_path := GetDBPath(db_name)
		fmt.Printf("数据库 %s 不存在\n", db_path)
		println("请重启本程序并选择 1 以初始化该数据库")
		WaitExit(1)
	}
}

func CheckDBInited(db *sql.DB, db_path string) {
	if tables, yes := DBInited(db); !yes {
		fmt.Printf("数据库 %s 缺少如下表：\n", db_path)
		fmt.Printf("%s\n", tables)
		println("请删除或备份该数据库文件，重启本程序并选择 1 以初始化该数据库")
		WaitExit(1)
	}
}

func CheckDBHasData(db *sql.DB, db_path string) {
	if tables, yes := DBHasData(db); !yes {
		fmt.Printf("数据库 %s 的如下表中还没有数据\n", db_path)
		fmt.Printf("%s\n", tables)
		println("请重启本程序并选择 2 以获取目录树")
		WaitExit(1)
	}
}

func CheckDBNoData(db *sql.DB, db_path string) {
	if tables, yes := DBNoData(db); !yes {
		fmt.Printf("数据库 %s 的如下表中存在数据\n", db_path)
		fmt.Printf("%s\n", tables)
		WaitExit(1)
	}
}

func CheckAllDBExists() {
	println("检查是否所有的数据库都存在")

	if paths, yes := AllDBExists(); !yes {
		println("检查到如下数据库还不存在：")
		fmt.Printf("%s\n", paths)
		println("请重启本程序并选择 1 以初始化数据库")
		WaitExit(1)
	}
}

func CheckAllDBInited() {
	println("检查是否所有的数据库都已初始化")

	if paths, yes := AllDBInited(); !yes {
		println("检查到如下数据库还没有初始化：")
		fmt.Printf("%s\n", paths)
		println("请重启本程序并选择 1 以初始化这些数据库")
		WaitExit(1)
	}
}

func CheckAllDBHasData() {
	println("检查所有的数据库的 dirs 表和 files 表中已有数据")

	if paths, yes := AllDBHasData(); !yes {
		println("检查到如下数据库的 dirs 或 files 表中还没有数据：")
		fmt.Printf("%s\n", paths)
		println("请重启本程序并选择 2 以获取目录树")
		WaitExit(1)
	}
}

func AllDBExists() ([]string, bool) {
	var paths []string = []string{}

	for db_name := range g_dbs {
		db_path := GetDBPath(db_name)

		if !DBExists(db_name) {
			paths = append(paths, db_path)
		}
	}

	return paths, len(paths) <= 0
}

func AllDBInited() ([]string, bool) {
	var paths []string = []string{}

	for db_name, db := range g_dbs {
		db_path := GetDBPath(db_name)

		if tables, yes := DBInited(db); !yes {
			paths = append(paths, db_path, strings.Join(tables, " "))
		}
	}

	return paths, len(paths) <= 0
}

func AllDBHasData() ([]string, bool) {
	var paths []string = []string{}

	for db_name, db := range g_dbs {
		db_path := GetDBPath(db_name)

		if tables, yes := DBHasData(db); !yes {
			paths = append(paths, db_path, strings.Join(tables, " "))
		}
	}

	return paths, len(paths) <= 0
}
