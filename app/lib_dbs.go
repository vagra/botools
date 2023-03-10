package app

import (
	"fmt"
	"os"
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

func GetOldDBPath(disk_name string) string {
	return fmt.Sprintf("%s/old.%s%s", DB_DIR, disk_name, DB_EXT)
}

func GetMemDBPath(disk_name string) string {
	return fmt.Sprintf("file:%s?mode=memory&cache=shared", disk_name)
}

func CheckTaskHasDBs() {
	if len(g_dbs) > 0 {
		return
	}

	println("没有需要处理的数据库")
	WaitExit(0)
}

func GetAllDBs() {
	g_dbs = make(map[string]*DB)

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

	g_dbs = make(map[string]*DB)

	for disk_name := range g_disks {

		db_path := GetDBPath(disk_name)

		if DBExists(disk_name) {
			fmt.Printf("数据库 %s 已存在，跳过\n", db_path)
			continue
		}

		fmt.Printf("数据库 %s 不存在，新建...\n", db_path)

		db := DBOpen(db_path)

		g_dbs[disk_name] = db
	}
}

func GetNotInitedDBs() {
	println("检查不存在或没有初始化的数据库")

	g_dbs = make(map[string]*DB)

	for disk_name := range g_disks {

		db_path := GetDBPath(disk_name)

		is_new := !DBExists(disk_name)

		db := DBOpen(db_path)

		if _, yes := db.Inited(); yes {
			continue
		}

		if is_new {
			fmt.Printf("数据库 %s 不存在，新建……\n", db_path)
		} else {
			fmt.Printf("数据库 %s 没有初始化，打开\n", db_path)
		}

		g_dbs[disk_name] = db
	}
}

func GetEmptyDBs() {
	println("获取所有 dirs 和 files 表为空的数据库")

	g_dbs = make(map[string]*DB)

	for disk_name := range g_disks {

		db_path := GetDBPath(disk_name)

		CheckDBExists(disk_name)

		db := DBOpen(db_path)

		CheckDBInited(db, db_path)

		if _, yes := db.NoData(); !yes {
			continue
		}
		println(db_path)

		g_dbs[disk_name] = db
	}
}

func GetInitedDBs() {
	println("获取所有已经初始化的数据库")

	g_dbs = make(map[string]*DB)

	for disk_name := range g_disks {

		db_path := GetDBPath(disk_name)

		CheckDBExists(disk_name)

		db := DBOpen(db_path)

		CheckDBInited(db, db_path)

		println(db_path)

		g_dbs[disk_name] = db
	}
}

func GetHasDataDBs() {
	println("获取所有 dirs 和 files 表都有数据的数据库")

	g_dbs = make(map[string]*DB)

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
	println("获取还有 files 没有 sha1 的数据库")

	g_dbs = make(map[string]*DB)

	for disk_name := range g_disks {

		db_path := GetDBPath(disk_name)

		CheckDBExists(disk_name)

		db := DBOpen(db_path)

		CheckDBInited(db, db_path)

		CheckDBHasData(db, db_path)

		if _, yes := db.NeedCheckSum(); !yes {
			continue
		}

		println(db_path)

		g_dbs[disk_name] = db
	}
}

func GetHasErrorDBs() {
	println("获取所有存在异常文件和文件夹的数据库")

	g_dbs = make(map[string]*DB)

	for disk_name := range g_errors {

		db_path := GetDBPath(disk_name)

		CheckDBExists(disk_name)

		db := DBOpen(db_path)

		CheckDBInited(db, db_path)

		CheckDBHasData(db, db_path)

		println(db_path)

		g_dbs[disk_name] = db
	}
}

func GetDBCounters() {
	println("获取每个数据库的 dirs 和 files 的 id 计数器")

	g_dirs_counter = make(map[string]*int64)
	g_files_counter = make(map[string]*int64)

	for disk_name, db := range g_dbs {

		dir_counter := db.QueryMaxDirIndex()
		file_counter := db.QueryMaxFileIndex()
		g_dirs_counter[disk_name] = &dir_counter
		g_files_counter[disk_name] = &file_counter

		fmt.Printf("%s index: dir %d, file %d\n", disk_name, dir_counter, file_counter)
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

func CheckOldDBExists(db_name string) {
	if !OldDBExists(db_name) {
		old_path := GetOldDBPath(db_name)
		fmt.Printf("备份数据库 %s 不存在\n", old_path)
		println("请重启本程序并选择 1 以初始化该数据库")
		WaitExit(1)
	}
}

func CheckDBInited(db *DB, db_path string) {
	if tables, yes := db.Inited(); !yes {
		fmt.Printf("数据库 %s 缺少如下表：\n", db_path)
		for _, table := range tables {
			println(table)
		}
		println("请删除或备份该数据库文件，重启本程序并选择 1 以初始化该数据库")
		WaitExit(1)
	}
}

func CheckDBHasData(db *DB, db_path string) {
	if tables, yes := db.HasData(); !yes {
		fmt.Printf("数据库 %s 的如下表还没有数据\n", db_path)
		for _, table := range tables {
			println(table)
		}
		println("请重启本程序并选择 2 以获取目录树")
		WaitExit(1)
	}
}

func CheckDBNoData(db *DB, db_path string) {
	if tables, yes := db.NoData(); !yes {
		fmt.Printf("数据库 %s 的如下表存在数据\n", db_path)
		for _, table := range tables {
			println(table)
		}
		WaitExit(1)
	}
}

func CheckAllDBExists() {
	println("检查是否所有数据库都存在")

	if paths, yes := AllDBExists(); !yes {
		println("检查到如下数据库还不存在：")
		for _, db_path := range paths {
			println(db_path)
		}
		println("请重启本程序并选择 1 以初始化数据库")
		WaitExit(1)
	}
}

func CheckAllDBInited() {
	println("检查是否所有数据库都已初始化")

	if paths, yes := AllDBInited(); !yes {
		println("检查到如下数据库还没有初始化：")
		for _, db_path := range paths {
			println(db_path)
		}
		println("请重启本程序并选择 1 以初始化这些数据库")
		WaitExit(1)
	}
}

func CheckAllDBHasData() {
	println("检查是否所有数据库的 dirs 表和 files 表都有数据")

	if paths, yes := AllDBHasData(); !yes {
		println("检查到如下数据库的 dirs 或 files 表还没有数据：")
		for _, db_path := range paths {
			println(db_path)
		}
		println("请重启本程序并选择 2 以获取目录树")
		WaitExit(1)
	}
}

func CheckAllDBRootPathCorrect() {
	fmt.Printf("检查是否所有数据库中的根路径都与 %s 中的一致\n", CONFIG_INI)

	if paths, yes := AllDBRootPathCorrect(); !yes {
		fmt.Printf("检查到如下数据库中的根路径与 %s 中的不一致：\n", CONFIG_INI)
		for db_path, root_path := range paths {
			fmt.Printf("%s\t%s\n", db_path, root_path)
		}
		fmt.Printf("请检查 %s 和以上数据库\n", CONFIG_INI)
		WaitExit(1)
	}
}

func AllDBExists() ([]string, bool) {
	var paths []string = []string{}

	for disk_name := range g_disks {
		db_path := GetDBPath(disk_name)

		if !DBExists(disk_name) {
			paths = append(paths, db_path)
		}
	}

	return paths, len(paths) <= 0
}

func AllDBInited() ([]string, bool) {
	var paths []string = []string{}

	for disk_name := range g_disks {
		db_path := GetDBPath(disk_name)

		db := DBOpen(db_path)

		if tables, yes := db.Inited(); !yes {
			info := fmt.Sprintf("%s\t%s", db_path, tables)
			paths = append(paths, info)
		}

		db.Close()
	}

	return paths, len(paths) <= 0
}

func AllDBHasData() ([]string, bool) {
	var paths []string = []string{}

	for disk_name := range g_disks {
		db_path := GetDBPath(disk_name)

		db := DBOpen(db_path)

		if tables, yes := db.HasData(); !yes {
			info := fmt.Sprintf("%s\t%s", db_path, tables)
			paths = append(paths, info)
		}

		db.Close()
	}

	return paths, len(paths) <= 0
}

func AllDBRootPathCorrect() (map[string]string, bool) {
	var paths map[string]string = make(map[string]string)

	for disk_name, disk_path := range g_disks {
		db_path := GetDBPath(disk_name)

		db := DBOpen(db_path)

		dir := db.GetRootDir()

		if dir.path != disk_path {
			paths[db_path] = dir.path
		}

		db.Close()
	}

	return paths, len(paths) <= 0
}

func GenDirUID(disk_name string) string {
	prefix := GetDirPrefix(disk_name)
	counter := GetDirsCounter(disk_name)
	return GenUID(prefix, counter)
}

func GenFileUID(disk_name string) string {
	prefix := GetFilePrefix(disk_name)
	counter := GetFilesCounter(disk_name)
	return GenUID(prefix, counter)
}

func GetDirUID(disk_name string, id int64) string {
	prefix := GetDirPrefix(disk_name)
	return GetUID(prefix, id)
}

func GetFileUID(disk_name string, id int64) string {
	prefix := GetFilePrefix(disk_name)
	return GetUID(prefix, id)
}

func GetDirPrefix(disk_name string) string {
	return fmt.Sprintf("%s-%s", disk_name, DIR_PRE)
}

func GetFilePrefix(disk_name string) string {
	return fmt.Sprintf("%s-%s", disk_name, FILE_PRE)
}

func GetDirsCounter(disk_name string) *int64 {
	return g_dirs_counter[disk_name]
}

func GetFilesCounter(disk_name string) *int64 {
	return g_files_counter[disk_name]
}

func GenUID(prefix string, counter *int64) string {
	*counter += 1
	return fmt.Sprintf("%s%08d", prefix, *counter)
}

func GetUID(prefix string, id int64) string {
	return fmt.Sprintf("%s%08d", prefix, id)
}
