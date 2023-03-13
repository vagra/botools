package app

import (
	"database/sql"
	"regexp"
)

type DB sql.DB

// --------------------------------------------
// query the dbs.
// --------------------------------------------

func DBExists(db_name string) bool {
	db_path := GetDBPath(db_name)

	return FileExists(db_path)
}

func OldDBExists(db_name string) bool {
	old_path := GetOldDBPath(db_name)

	return FileExists(old_path)
}

func DBOpen(db_path string) *DB {
	db, err := sql.Open("sqlite3", db_path)
	Check(err, "open db %s failed", db_path)

	return (*DB)(db)
}

func (db *DB) Close() error {
	err := (*sql.DB)(db).Close()

	return err
}

func (db *DB) Inited() ([]string, bool) {
	var tables []string = []string{}

	if !db.TableExists("dirs") {
		tables = append(tables, "dirs")
	}
	if !db.TableExists("files") {
		tables = append(tables, "files")
	}
	if !db.TableExists("infos") {
		tables = append(tables, "infos")
	}

	return tables, len(tables) <= 0
}

func (db *DB) HasData() ([]string, bool) {
	var tables []string = []string{}

	if db.QueryDirsCount() <= 0 {
		tables = append(tables, "dirs")
	}

	if db.QueryFilesCount() <= 0 {
		tables = append(tables, "files")
	}

	return tables, len(tables) <= 0
}

func (db *DB) NoData() ([]string, bool) {
	var tables []string = []string{}

	if db.QueryDirsCount() > 0 {
		tables = append(tables, "dirs")
	}

	if db.QueryFilesCount() > 0 {
		tables = append(tables, "files")
	}

	return tables, len(tables) <= 0
}

//////// files

func (db *DB) NeedCheckSum() (int64, bool) {

	count := db.QueryNoSHA1FilesCount()

	return count, count > 0
}

// --------------------------------------------
// query if the tables exists.
// --------------------------------------------

func (db *DB) TableExists(table_name string) bool {
	var name string

	row := db.QueryRow(SQL_CHECK_TABLE, table_name)

	err := row.Scan(&name)

	return err == nil
}

func (db *DB) DirsTableExists() bool {
	return db.TableExists("dirs")
}

func (db *DB) FilesTableExists() bool {
	return db.TableExists("files")
}

func (db *DB) InfosTableExists() bool {
	return db.TableExists("infos")
}

// --------------------------------------------
// create the tables.
// --------------------------------------------

func (db *DB) CreateDirsTable() {
	db.Exec(SQL_CREATE_DIRS)
}

func (db *DB) CreateFilesTable() {
	db.Exec(SQL_CREATE_FILES)
}

func (db *DB) CreateVDirsTable() {
	db.Exec(SQL_CREATE_VDIRS)
}

func (db *DB) CreateVFilesTable() {
	db.Exec(SQL_CREATE_VFILES)
}

func (db *DB) CreateInfosTable() {
	db.Exec(SQL_CREATE_INFOS)
}

// --------------------------------------------
// query the tables.
// --------------------------------------------

//////// dirs

func (db *DB) GetRootDir() (*Dir, bool) {
	var dir Dir

	row := db.QueryRow(SQL_GET_ROOT_DIR)
	ok := DBScanRow(row, SQL_GET_ROOT_DIR,
		&dir.id, &dir.parent_id, &dir.name, &dir.path, &dir.status, &dir.error)

	return &dir, ok
}

func (db *DB) QueryDirsCount() int64 {
	var count int64 = 0

	row := db.QueryRow(SQL_COUNT_DIRS)
	DBScanRow(row, SQL_COUNT_DIRS, &count)

	return count
}

func (db *DB) QueryMaxDirIndex() int64 {
	var id string = ""

	row := db.QueryRow(SQL_MAX_DIR_ID)
	DBScanRow(row, SQL_MAX_DIR_ID, &id)

	index, _ := DBGetIDIndex(id)

	return index
}

func (db *DB) GetAllDirs() map[string]*Dir {

	var dirs map[string]*Dir = make(map[string]*Dir)

	rows := db.QueryRows(SQL_GET_ALL_DIRS)
	defer rows.Close()

	for rows.Next() {
		var dir Dir

		DBScanRows(rows, SQL_GET_ALL_DIRS,
			&dir.id, &dir.parent_id, &dir.name, &dir.path, &dir.status, &dir.error)

		dirs[dir.id] = &dir
	}

	return dirs
}

func (db *DB) QueryDirIDFromPath(path string) (string, bool) {
	var id string = ""

	row := db.QueryRow(SQL_PATH_GET_DIR_ID, path)
	DBScanRow(row, SQL_PATH_GET_DIR_ID, &id)

	return id, len(id) > 8
}

func (db *DB) GetADirID() string {
	var id string = ""

	row := db.QueryRow(SQL_GET_A_DIR_ID)
	DBScanRow(row, SQL_GET_A_DIR_ID, &id)

	return id
}

func (db *DB) GetDirIDPrefix() string {
	id := db.GetADirID()

	prefix, _ := DBGetIDPrefix(id)

	return prefix
}

func (db *DB) GetNextDir(id string) (*Dir, bool) {
	var dir Dir

	row := db.QueryRow(SQL_GET_NEXT_DIR, id)
	ok := DBScanRow(row, SQL_GET_NEXT_DIR,
		&dir.id, &dir.parent_id, &dir.name, &dir.path, &dir.status, &dir.error)

	return &dir, ok
}

//////// files

func (db *DB) QueryFilesCount() int64 {
	var count int64 = 0

	row := db.QueryRow(SQL_COUNT_FILES)
	DBScanRow(row, SQL_COUNT_FILES, &count)

	return count
}

func (db *DB) QueryMaxFileIndex() int64 {
	var id string = ""

	row := db.QueryRow(SQL_MAX_FILE_ID)
	DBScanRow(row, SQL_MAX_FILE_ID, &id)

	index, _ := DBGetIDIndex(id)

	return index
}

func (db *DB) GetAllFiles() map[string]*File {

	var files map[string]*File = make(map[string]*File)

	rows := db.QueryRows(SQL_GET_ALL_FILES)
	defer rows.Close()

	for rows.Next() {
		var file File

		DBScanRows(rows, SQL_GET_ALL_FILES,
			&file.id, &file.parent_id, &file.name, &file.path, &file.status, &file.error)

		files[file.id] = &file
	}

	return files
}

func (db *DB) QueryNoSHA1FilesCount() int64 {
	var count int64 = 0

	row := db.QueryRow(SQL_GET_NO_SHA1_FILES_COUNT)
	DBScanRow(row, SQL_GET_NO_SHA1_FILES_COUNT, &count)

	return count
}

func (db *DB) GetNoSHA1Files() map[string]*File {

	var files map[string]*File = make(map[string]*File)

	rows := db.QueryRows(SQL_GET_NO_SHA1_FILES)
	defer rows.Close()

	for rows.Next() {
		var file File

		DBScanRows(rows, SQL_GET_NO_SHA1_FILES,
			&file.id, &file.parent_id, &file.name, &file.path)

		files[file.id] = &file
	}

	return files
}

func (db *DB) QueryFileIDFromPath(path string) (string, bool) {
	var id string = ""

	row := db.QueryRow(SQL_PATH_GET_FILE_ID, path)
	DBScanRow(row, SQL_PATH_GET_FILE_ID, &id)

	return id, len(id) > 8
}

func (db *DB) GetAFileID() string {
	var id string = ""

	row := db.QueryRow(SQL_GET_A_FILE_ID)
	DBScanRow(row, SQL_GET_A_FILE_ID, &id)

	return id
}

func (db *DB) GetFileIDPrefix() string {
	id := db.GetAFileID()

	prefix, _ := DBGetIDPrefix(id)

	return prefix
}

func (db *DB) GetNextFile(id string) (*File, bool) {
	var file File

	row := db.QueryRow(SQL_GET_NEXT_FILE, id)

	ok := DBScanRow(row, SQL_GET_NEXT_FILE,
		&file.id, &file.parent_id, &file.name, &file.path,
		&file.size, &file.status, &file.error, &file.dup_id, &file.sha1)

	return &file, ok
}

func (db *DB) GetNextNodupFile(id string) (*File, bool) {
	var file File

	row := db.QueryRow(SQL_GET_NEXT_NODUP_FILE, id)

	ok := DBScanRow(row, SQL_GET_NEXT_NODUP_FILE,
		&file.id, &file.parent_id, &file.name, &file.path,
		&file.size, &file.status, &file.error, &file.dup_id, &file.sha1)

	return &file, ok
}

func (db *DB) GetNextDupFile(id string) (*File, bool) {
	var file File

	row := db.QueryRow(SQL_GET_NEXT_DUP_FILE, id)

	ok := DBScanRow(row, SQL_GET_NEXT_DUP_FILE,
		&file.id, &file.parent_id, &file.name, &file.path,
		&file.size, &file.status, &file.error, &file.dup_id, &file.sha1)

	return &file, ok
}

//////// infos

func (db *DB) GetVersion() int {
	var version int

	row := db.QueryRow(SQL_GET_VERSION)
	DBScanRow(row, SQL_GET_VERSION, &version)

	return version
}

// --------------------------------------------
// insert into tables.
// --------------------------------------------

//////// dirs

func (db *DB) AddDir(dir *Dir) {
	db.Exec(SQL_ADD_DIR,
		dir.id, dir.parent_id, dir.name, dir.path, dir.mod_time)
}

//////// files

func (db *DB) AddFile(file *File) {
	db.Exec(SQL_ADD_FILE,
		file.id, file.parent_id, file.name, file.path, file.size, file.mod_time)
}

//////// infos

func (db *DB) AddInfo(version int) {
	db.Exec(SQL_ADD_INFO, version)
}

// --------------------------------------------
// update the tables.
// --------------------------------------------

//////// dirs

func (db *DB) ModRootDir(path string) {
	db.Exec(SQL_MOD_ROOT_DIR, path, path)
}

func (db *DB) TrimDirsID() {
	db.Exec(SQL_TRIM_DIRS_ID)
}

func (db *DB) ModDirsStatus(status int) {
	db.Exec(SQL_MOD_DIRS_STATUS, status)
}

func (db *DB) ReplaceDirsPath(src string, dst string) {
	db.Exec(SQL_REPLACE_DIRS_PATH, dst, src, src)
}

func (db *DB) ModDirStatus(id string, status int8) {
	db.Exec(SQL_MOD_DIR_STATUS, status, id)
}

func (db *DB) ModDirError(id string, code int) {
	db.Exec(SQL_MOD_DIR_ERROR, code, id)
}

func (db *DB) ModDirsDiskID(src string, dst string) {
	db.Exec(SQL_REPLACE_DIRS_ID, src, dst)
	db.Exec(SQL_REPLACE_DIRS_PARENT_ID, src, dst)

	db.Exec(SQL_REPLACE_FILES_PARENT_ID, src, dst)
}

//////// files

func (db *DB) TrimFilesID() {
	db.Exec(SQL_TRIM_FILES_ID)
}

func (db *DB) ModFilesStatus(status int) {
	db.Exec(SQL_MOD_FILES_STATUS, status)
}

func (db *DB) ReplaceFilesPath(src string, dst string) {
	db.Exec(SQL_REPLACE_FILES_PATH, dst, src, src)
}

func (db *DB) ModFilesDiskID(src string, dst string) {
	db.Exec(SQL_REPLACE_FILES_ID, src, dst)
}

func (db *DB) ModDirFilesError(id string, code int8) {
	db.Exec(SQL_MOD_DIR_FILES_ERROR, code, id)
}

func (db *DB) ModFileSha1(id string, sha1 string) {
	db.Exec(SQL_MOD_FILE_SHA1, sha1, id)
}

func (db *DB) ModFileStatus(id string, status int8) {
	db.Exec(SQL_MOD_FILE_STATUS, status, id)
}

func (db *DB) ModFileError(id string, code int8) {
	db.Exec(SQL_MOD_FILE_ERROR, code, id)
}

func (db *DB) ModFileDupID(id string, dup_id string) {
	db.Exec(SQL_MOD_FILE_DUP_ID, dup_id, id)
}

func (db *DB) BulkModFilesSha1(files *[]*File) {
	db.BeginBulk()

	for _, file := range *files {
		db.ModFileSha1OrStatus(file)
	}

	db.EndBulk()
}

func (db *DB) ModFileSha1OrStatus(file *File) {
	if file.status == 0 {
		db.ModFileSha1(file.id, file.sha1)
	} else {
		db.ModFileStatus(file.id, file.status)
	}
}

//////// infos

func (db *DB) ModVersion(version int) {
	db.Exec(SQL_MOD_VERSION, version)
}

// --------------------------------------------
// common function
// --------------------------------------------

func (db *DB) BeginBulk() {
	db.Exec(SQL_BEGIN)
}

func (db *DB) EndBulk() {
	db.Exec(SQL_END)
}

func (db *DB) Exec(sql_name string, args ...interface{}) {
	_, err := g_dot.Exec((*sql.DB)(db), sql_name, args...)
	Check(err, "db error when run SQL %s", sql_name)
}

func (db *DB) QueryRow(sql_name string, args ...interface{}) *sql.Row {
	row, err := g_dot.QueryRow((*sql.DB)(db), sql_name, args...)
	Check(err, "db error when run SQL %s", sql_name)

	return row
}

func DBScanRow(row *sql.Row, sql_name string, dest ...any) bool {
	err := row.Scan(dest...)
	return err == nil
}

func (db *DB) QueryRows(sql_name string, args ...interface{}) *sql.Rows {
	rows, err := g_dot.Query((*sql.DB)(db), sql_name, args...)
	Check(err, "db error when run SQL %s", sql_name)

	return rows
}

func DBScanRows(rows *sql.Rows, sql_name string, dest ...any) bool {
	err := rows.Scan(dest...)
	return err == nil
}

func DBGetIDPrefix(id string) (string, bool) {
	regex := regexp.MustCompile(ID_REGEX)
	matches := regex.FindStringSubmatch(id)

	if len(matches) < 3 {
		return "ERROR NO MATCHES", false
	}

	return matches[1], true
}

func DBGetIDSuffix(id string) (string, bool) {
	regex := regexp.MustCompile(ID_REGEX)
	matches := regex.FindStringSubmatch(id)

	if len(matches) < 3 {
		return "ERROR NO MATCHES", false
	}

	return matches[2], true
}

func DBGetIDIndex(id string) (int64, bool) {
	code, ok := DBGetIDSuffix(id)
	if !ok {
		return 0, false
	}

	return Str2Int64(code), true
}
