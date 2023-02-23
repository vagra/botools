package app

import (
	"database/sql"
	"fmt"
	"log"
)

type Disk struct {
	id   int64
	name string
	path string
	size int64
}

type Info struct {
	id         int64
	db_version int
}

type Dir struct {
	id        string
	parent_id string
	name      string
	path      string
	size      int64
	status    int8
	mod_time  string
}

type File struct {
	id        string
	parent_id string
	name      string
	path      string
	size      int64
	status    int8
	sha1      string
	mod_time  string
}

func (d Dir) Tuple() string {
	return fmt.Sprintf("('%s', '%s', '%s', '%s', '%d', '%d', '%s')",
		d.id, d.parent_id, d.name, d.path, d.size, d.status, d.mod_time)
}

func (d Dir) AddMarks(marks *[]string) {
	*marks = append(*marks, "(?, ?, ?, ?, ?)")
}

func (d Dir) AddArgs(args *[]interface{}) {
	*args = append(*args, d.id)
	*args = append(*args, d.parent_id)
	*args = append(*args, d.name)
	*args = append(*args, d.path)
	*args = append(*args, d.mod_time)
}

func (f File) Tuple() string {
	return fmt.Sprintf("('%s', '%s', '%s', '%s', '%d', '%s', '%s')",
		f.id, f.parent_id, f.name, f.path, f.size, f.sha1, f.mod_time)
}

func (f File) AddMarks(marks *[]string) {
	*marks = append(*marks, "(?, ?, ?, ?, ?, ?)")
}

func (f File) AddArgs(args *[]interface{}) {
	*args = append(*args, f.id)
	*args = append(*args, f.parent_id)
	*args = append(*args, f.name)
	*args = append(*args, f.path)
	*args = append(*args, f.size)
	*args = append(*args, f.mod_time)
}

func DBTableExists(db *sql.DB, table_name string) bool {
	var name string

	row, err := g_dot.QueryRow(db, SQL_CHECK_TABLE, table_name)
	Check(err, "执行 SQL "+SQL_CHECK_TABLE+" 时出错")

	err = row.Scan(&name)
	if err != nil {
		return false
	}

	return true
}

func DBInfosTableExists(db *sql.DB) bool {
	return DBTableExists(db, "infos")
}

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

func DBQueryDirsCount(db *sql.DB) int64 {

	var count int64 = 0

	row, err := g_dot.QueryRow(db, SQL_COUNT_DIRS)
	Check(err, "执行 SQL "+SQL_COUNT_DIRS+" 时出错")

	err = row.Scan(&count)
	Check(err, "执行 SQL "+SQL_COUNT_DIRS+" 后获取 count 时出错")

	return count
}

func DBQueryFilesCount(db *sql.DB) int64 {

	var count int64 = 0

	row, err := g_dot.QueryRow(db, SQL_COUNT_FILES)
	Check(err, "执行 SQL "+SQL_COUNT_FILES+" 时出错")

	err = row.Scan(&count)
	Check(err, "执行 SQL "+SQL_COUNT_FILES+" 后获取 count 时出错")

	return count
}

func DBUpdateFile(db *sql.DB, file *File) {
	if file.status == 0 {
		DBUpdateFileSha1(db, file)
	} else {
		DBUpdateFileStatus(db, file)
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

func DBTrimDirIDs(db *sql.DB) {
	_, err := g_dot.Exec(db, SQL_TRIM_DIR_IDS)
	if err != nil {
		log.Printf("db trim dir ids error: %s\n" + err.Error())
	}
}

func DBTrimFileIDs(db *sql.DB) {
	_, err := g_dot.Exec(db, SQL_TRIM_FILE_IDS)
	if err != nil {
		log.Printf("db trim file ids error: %s\n" + err.Error())
	}
}

func DBReplaceDirPaths(db *sql.DB, src string, dst string) {
	_, err := g_dot.Exec(db, SQL_REPLACE_DIR_PATHS, src, dst)
	if err != nil {
		log.Printf("db replace dir paths error: %s\n" + err.Error())
	}
}

func DBReplaceFilePaths(db *sql.DB, src string, dst string) {
	_, err := g_dot.Exec(db, SQL_REPLACE_FILE_PATHS, src, dst)
	if err != nil {
		log.Printf("db replace file paths error: %s\n" + err.Error())
	}
}

func DBGetRootDir(db *sql.DB) Dir {
	row, err := g_dot.QueryRow(db, SQL_GET_ROOT_DIR)
	Check(err, "执行 SQL "+SQL_GET_ROOT_DIR+" 时出错")

	var dir Dir

	err = row.Scan(&dir.id, &dir.parent_id, &dir.name, &dir.path)
	Check(err, "执行 SQL "+SQL_GET_ROOT_DIR+" 后获取 dir 时出错")

	return dir
}

func DBUpdateRootDir(db *sql.DB, path string) {
	_, err := g_dot.Exec(db, SQL_MOD_ROOT_DIR, path, path)
	if err != nil {
		log.Printf("db update root dir error: %s\n" + err.Error())
	}
}

func DBGetAllDirs(db *sql.DB) map[string]*Dir {

	var dirs map[string]*Dir = make(map[string]*Dir)

	rows, err := g_dot.Query(db, SQL_GET_ALL_DIRS)
	Check(err, "执行 SQL "+SQL_GET_ALL_DIRS+" 时出错")
	defer rows.Close()

	for rows.Next() {
		var dir Dir

		err = rows.Scan(&dir.id, &dir.parent_id, &dir.name, &dir.path)
		Check(err, "执行 SQL "+SQL_GET_ALL_DIRS+" 后获取 dir 时出错")

		dirs[dir.id] = &dir
	}

	return dirs
}

func DBGetFilesNoSHA1(db *sql.DB) map[string]*File {

	var files map[string]*File = make(map[string]*File)

	rows, err := g_dot.Query(db, SQL_GET_FILES_NO_SHA1)
	Check(err, "执行 SQL "+SQL_GET_FILES_NO_SHA1+" 时出错")
	defer rows.Close()

	for rows.Next() {
		var file File

		err = rows.Scan(&file.id, &file.parent_id, &file.name, &file.path)
		Check(err, "执行 SQL "+SQL_GET_FILES_NO_SHA1+" 后获取 file 时出错")

		files[file.id] = &file
	}

	return files
}

func DBAddInfo(db *sql.DB, version int) {
	_, err := g_dot.Exec(db, SQL_ADD_INFO, version)
	if err != nil {
		log.Printf("db add info to infos table error: %s\n" + err.Error())
	}
}

func DBGetVersion(db *sql.DB) int {
	row, err := g_dot.QueryRow(db, SQL_GET_VERSION)
	Check(err, "执行 SQL "+SQL_GET_VERSION+" 时出错")

	var version int

	err = row.Scan(&version)
	Check(err, "执行 SQL "+SQL_GET_VERSION+" 后获取 db_version 时出错")

	return version
}

func DBUpdateVersion(db *sql.DB, version int) {
	_, err := g_dot.Exec(db, SQL_MOD_VERSION, version)
	if err != nil {
		log.Printf("db update infos.db_version error: %s\n" + err.Error())
	}
}
