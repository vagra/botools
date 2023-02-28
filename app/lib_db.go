package app

import (
	"database/sql"
	"log"
)

//// query the dbs.

func DBExists(db_name string) bool {
	db_path := GetDBPath(db_name)

	if !FileExists(db_path) {
		return false
	}

	return true
}

func DBOpen(db_path string) *sql.DB {
	db, err := sql.Open("sqlite3", db_path)
	Check(err, "打开数据库 %s 失败", db_path)

	return db
}

func DBInited(db *sql.DB) ([]string, bool) {
	var tables []string = []string{}

	if !DBTableExists(db, "dirs") {
		tables = append(tables, "dirs")
	}
	if !DBTableExists(db, "files") {
		tables = append(tables, "files")
	}
	if !DBTableExists(db, "infos") {
		tables = append(tables, "infos")
	}

	return tables, len(tables) <= 0
}

func DBHasData(db *sql.DB) ([]string, bool) {
	var tables []string = []string{}

	if DBQueryDirsCount(db) <= 0 {
		tables = append(tables, "dirs")
	}

	if DBQueryFilesCount(db) <= 0 {
		tables = append(tables, "files")
	}

	return tables, len(tables) <= 0
}

func DBNoData(db *sql.DB) ([]string, bool) {
	var tables []string = []string{}

	if DBQueryDirsCount(db) > 0 {
		tables = append(tables, "dirs")
	}

	if DBQueryFilesCount(db) > 0 {
		tables = append(tables, "files")
	}

	return tables, len(tables) <= 0
}

// files

func DBNeedCheckSum(db *sql.DB) (int64, bool) {

	count := DBQueryNoSHA1FilesCount(db)

	return count, count > 0
}

//// query if the tables exists.

func DBTableExists(db *sql.DB, table_name string) bool {
	var name string

	row, err := g_dot.QueryRow(db, SQL_CHECK_TABLE, table_name)
	Check(err, "执行 SQL %s 时出错", SQL_CHECK_TABLE)

	err = row.Scan(&name)
	if err != nil {
		return false
	}

	return true
}

func DBDirsTableExists(db *sql.DB) bool {
	return DBTableExists(db, "dirs")
}

func DBFilesTableExists(db *sql.DB) bool {
	return DBTableExists(db, "files")
}

func DBInfosTableExists(db *sql.DB) bool {
	return DBTableExists(db, "infos")
}

//// create the tables.

func DBCreateDirsTable(db *sql.DB) {
	_, err := g_dot.Exec(db, SQL_CREATE_DIRS)
	Check(err, "创建 dirs 表失败")
}

func DBCreateFilesTable(db *sql.DB) {
	_, err := g_dot.Exec(db, SQL_CREATE_FILES)
	Check(err, "创建 files 表失败")
}

func DBCreateInfosTable(db *sql.DB) {
	_, err := g_dot.Exec(db, SQL_CREATE_INFOS)
	Check(err, "创建 infos 表失败")
}

//// query the tables.

// dirs

func DBGetRootDir(db *sql.DB) Dir {
	row, err := g_dot.QueryRow(db, SQL_GET_ROOT_DIR)
	Check(err, "执行 SQL %s 时出错", SQL_GET_ROOT_DIR)

	var dir Dir

	err = row.Scan(&dir.id, &dir.parent_id, &dir.name, &dir.path)
	Check(err, "执行 SQL %s 后获取 dir 时出错", SQL_GET_ROOT_DIR)

	return dir
}

func DBQueryDirsCount(db *sql.DB) int64 {

	var count int64 = 0

	row, err := g_dot.QueryRow(db, SQL_COUNT_DIRS)
	Check(err, "执行 SQL %s 时出错", SQL_COUNT_DIRS)

	err = row.Scan(&count)
	Check(err, "执行 SQL %s 后获取 count 时出错", SQL_COUNT_DIRS)

	return count
}

func DBGetAllDirs(db *sql.DB) map[string]*Dir {

	var dirs map[string]*Dir = make(map[string]*Dir)

	rows, err := g_dot.Query(db, SQL_GET_ALL_DIRS)
	Check(err, "执行 SQL %s 时出错", SQL_GET_ALL_DIRS)
	defer rows.Close()

	for rows.Next() {
		var dir Dir

		err = rows.Scan(&dir.id, &dir.parent_id, &dir.name, &dir.path)
		Check(err, "执行 SQL %s 后获取 dir 时出错", SQL_GET_ALL_DIRS)

		dirs[dir.id] = &dir
	}

	return dirs
}

// files

func DBQueryFilesCount(db *sql.DB) int64 {

	var count int64 = 0

	row, err := g_dot.QueryRow(db, SQL_COUNT_FILES)
	Check(err, "执行 SQL %s 时出错", SQL_COUNT_FILES)

	err = row.Scan(&count)
	Check(err, "执行 SQL %s 后获取 count 时出错", SQL_COUNT_FILES)

	return count
}

func DBQueryNoSHA1FilesCount(db *sql.DB) int64 {
	var count int64 = 0

	row, err := g_dot.QueryRow(db, SQL_GET_NO_SHA1_FILES_COUNT)
	Check(err, "执行 SQL %s 时出错", SQL_GET_NO_SHA1_FILES_COUNT)

	err = row.Scan(&count)
	Check(err, "执行 SQL %s 后获取 count 时出错", SQL_COUNT_FILES)

	return count
}

func DBGetNoSHA1Files(db *sql.DB) map[string]*File {

	var files map[string]*File = make(map[string]*File)

	rows, err := g_dot.Query(db, SQL_GET_NO_SHA1_FILES)
	Check(err, "执行 SQL %s 时出错", SQL_GET_NO_SHA1_FILES)
	defer rows.Close()

	for rows.Next() {
		var file File

		err = rows.Scan(&file.id, &file.parent_id, &file.name, &file.path)
		Check(err, "执行 SQL %s 后获取 file 时出错", SQL_GET_NO_SHA1_FILES)

		files[file.id] = &file
	}

	return files
}

// infos

func DBGetVersion(db *sql.DB) int {
	row, err := g_dot.QueryRow(db, SQL_GET_VERSION)
	Check(err, "执行 SQL %s 时出错", SQL_GET_VERSION)

	var version int

	err = row.Scan(&version)
	Check(err, "执行 SQL %s 后获取 db_version 时出错", SQL_GET_VERSION)

	return version
}

//// insert into tables.

// infos

func DBAddInfo(db *sql.DB, version int) {
	_, err := g_dot.Exec(db, SQL_ADD_INFO, version)
	if err != nil {
		log.Printf("db add info to infos table error: %s\n" + err.Error())
	}
}

//// update the tables.

// dirs

func DBUpdateRootDir(db *sql.DB, path string) {
	_, err := g_dot.Exec(db, SQL_MOD_ROOT_DIR, path, path)
	if err != nil {
		log.Printf("db update root dir error: %s\n" + err.Error())
	}
}

func DBTrimDirIDs(db *sql.DB) {
	_, err := g_dot.Exec(db, SQL_TRIM_DIR_IDS)
	if err != nil {
		log.Printf("db trim dir ids error: %s\n" + err.Error())
	}
}

func DBReplaceDirPaths(db *sql.DB, src string, dst string) {
	_, err := g_dot.Exec(db, SQL_REPLACE_DIR_PATHS, src, dst)
	if err != nil {
		log.Printf("db replace dir paths error: %s\n" + err.Error())
	}
}

// files

func DBUpdateFile(db *sql.DB, file *File) {
	if file.status == 0 {
		DBUpdateFileSha1(db, file)
	} else {
		DBUpdateFileStatus(db, file)
	}
}

func DBTrimFileIDs(db *sql.DB) {
	_, err := g_dot.Exec(db, SQL_TRIM_FILE_IDS)
	if err != nil {
		log.Printf("db trim file ids error: %s\n" + err.Error())
	}
}

func DBReplaceFilePaths(db *sql.DB, src string, dst string) {
	_, err := g_dot.Exec(db, SQL_REPLACE_FILE_PATHS, src, dst)
	if err != nil {
		log.Printf("db replace file paths error: %s\n" + err.Error())
	}
}

func DBUpdateFileSha1(db *sql.DB, file *File) {
	_, err := g_dot.Exec(db, SQL_MOD_FILE_SHA1, file.sha1, file.id)
	if err != nil {
		log.Printf("db update file sha1 error: %s\n", err.Error())
	}
}

func DBUpdateFileStatus(db *sql.DB, file *File) {
	_, err := g_dot.Exec(db, SQL_MOD_FILE_STATUS, file.status, file.id)
	if err != nil {
		log.Printf("db update file status error: %s\n" + err.Error())
	}
}

// infos

func DBUpdateVersion(db *sql.DB, version int) {
	_, err := g_dot.Exec(db, SQL_MOD_VERSION, version)
	if err != nil {
		log.Printf("db update infos.db_version error: %s\n" + err.Error())
	}
}
