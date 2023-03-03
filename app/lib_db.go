package app

import (
	"database/sql"
)

//// query the dbs.

func DBExists(db_name string) bool {
	db_path := GetDBPath(db_name)

	return FileExists(db_path)
}

func DBOpen(db_path string) *sql.DB {
	db, err := sql.Open("sqlite3", db_path)
	Check(err, "open db %s failed", db_path)

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

	row := DBQueryRow(db, SQL_CHECK_TABLE, table_name)

	err := row.Scan(&name)

	return err == nil
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
	DBExec(db, SQL_CREATE_DIRS)
}

func DBCreateFilesTable(db *sql.DB) {
	DBExec(db, SQL_CREATE_FILES)
}

func DBCreateVDirsTable(db *sql.DB) {
	DBExec(db, SQL_CREATE_VDIRS)
}

func DBCreateVFilesTable(db *sql.DB) {
	DBExec(db, SQL_CREATE_VFILES)
}

func DBCreateInfosTable(db *sql.DB) {
	DBExec(db, SQL_CREATE_INFOS)
}

//// query the tables.

// dirs

func DBGetRootDir(db *sql.DB) Dir {
	var dir Dir

	row := DBQueryRow(db, SQL_GET_ROOT_DIR)
	DBScanRow(row, SQL_GET_ROOT_DIR,
		&dir.id, &dir.parent_id, &dir.name, &dir.path)

	return dir
}

func DBQueryDirsCount(db *sql.DB) int64 {
	var count int64 = 0

	row := DBQueryRow(db, SQL_COUNT_DIRS)
	DBScanRow(row, SQL_COUNT_DIRS, &count)

	return count
}

func DBGetAllDirs(db *sql.DB) map[string]*Dir {

	var dirs map[string]*Dir = make(map[string]*Dir)

	rows := DBQueryRows(db, SQL_GET_ALL_DIRS)
	defer rows.Close()

	for rows.Next() {
		var dir Dir

		DBScanRows(rows, SQL_GET_ALL_DIRS,
			&dir.id, &dir.parent_id, &dir.name, &dir.path)

		dirs[dir.id] = &dir
	}

	return dirs
}

func DBQueryDirIDFromPath(db *sql.DB) string {
	var id string = ""

	row := DBQueryRow(db, SQL_PATH_GET_DIR_ID)
	DBScanRow(row, SQL_PATH_GET_DIR_ID, &id)

	return id
}

// files

func DBQueryFilesCount(db *sql.DB) int64 {
	var count int64 = 0

	row := DBQueryRow(db, SQL_COUNT_FILES)
	DBScanRow(row, SQL_COUNT_FILES, &count)

	return count
}

func DBQueryNoSHA1FilesCount(db *sql.DB) int64 {
	var count int64 = 0

	row := DBQueryRow(db, SQL_GET_NO_SHA1_FILES_COUNT)
	DBScanRow(row, SQL_GET_NO_SHA1_FILES_COUNT, &count)

	return count
}

func DBGetNoSHA1Files(db *sql.DB) map[string]*File {

	var files map[string]*File = make(map[string]*File)

	rows := DBQueryRows(db, SQL_GET_NO_SHA1_FILES)
	defer rows.Close()

	for rows.Next() {
		var file File

		DBScanRows(rows, SQL_GET_NO_SHA1_FILES,
			&file.id, &file.parent_id, &file.name, &file.path)

		files[file.id] = &file
	}

	return files
}

func DBQueryFileIDFromPath(db *sql.DB) string {
	var id string = ""

	row := DBQueryRow(db, SQL_PATH_GET_FILE_ID)
	DBScanRow(row, SQL_PATH_GET_FILE_ID, &id)

	return id
}

// infos

func DBGetVersion(db *sql.DB) int {
	var version int

	row := DBQueryRow(db, SQL_GET_VERSION)
	DBScanRow(row, SQL_GET_VERSION, &version)

	return version
}

//// insert into tables.

// infos

func DBAddInfo(db *sql.DB, version int) {
	DBExec(db, SQL_ADD_INFO, version)
}

//// update the tables.

// dirs

func DBUpdateRootDir(db *sql.DB, path string) {
	DBExec(db, SQL_MOD_ROOT_DIR, path, path)
}

func DBTrimDirIDs(db *sql.DB) {
	DBExec(db, SQL_TRIM_DIR_IDS)
}

func DBReplaceDirPaths(db *sql.DB, src string, dst string) {
	DBExec(db, SQL_REPLACE_DIR_PATHS, src, dst)
}

func DBUpdateDirError(db *sql.DB, dir *Dir, code int) {
	DBExec(db, SQL_MOD_DIR_ERROR, code, dir.id)
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
	DBExec(db, SQL_TRIM_FILE_IDS)
}

func DBReplaceFilePaths(db *sql.DB, src string, dst string) {
	DBExec(db, SQL_REPLACE_FILE_PATHS, src, dst)
}

func DBUpdateFileSha1(db *sql.DB, file *File) {
	DBExec(db, SQL_MOD_FILE_SHA1, file.sha1, file.id)
}

func DBUpdateFileStatus(db *sql.DB, file *File) {
	DBExec(db, SQL_MOD_FILE_STATUS, file.status, file.id)
}

func DBUpdateFileError(db *sql.DB, file *File, code int) {
	DBExec(db, SQL_MOD_FILE_ERROR, code, file.id)
}

func DBUpdateFileDupID(db *sql.DB, file *File, dup_id string) {
	DBExec(db, SQL_MOD_FILE_ERROR, dup_id, file.id)
}

// infos

func DBUpdateVersion(db *sql.DB, version int) {
	DBExec(db, SQL_MOD_VERSION, version)
}

// // common

func DBExec(db *sql.DB, sql_name string, args ...interface{}) {
	_, err := g_dot.Exec(db, sql_name, args...)
	Check(err, "db error when run SQL %s", sql_name)
}

func DBQueryRow(db *sql.DB, sql_name string, args ...interface{}) *sql.Row {
	row, err := g_dot.QueryRow(db, sql_name, args...)
	Check(err, "db error when run SQL %s", sql_name)

	return row
}

func DBScanRow(row *sql.Row, sql_name string, dest ...any) {
	err := row.Scan(dest...)
	Check(err, "db fetch result error after run SQL %s", sql_name)
}

func DBQueryRows(db *sql.DB, sql_name string, args ...interface{}) *sql.Rows {
	rows, err := g_dot.Query(db, sql_name, args...)
	Check(err, "db error when run SQL %s", sql_name)

	return rows
}

func DBScanRows(rows *sql.Rows, sql_name string, dest ...any) {
	err := rows.Scan(dest...)
	Check(err, "db fetch result error after run SQL %s", sql_name)
}
