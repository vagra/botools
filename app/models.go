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

type Dir struct {
	id        string
	parent_id string
	name      string
	path      string
	size      int64
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
	return fmt.Sprintf("('%s', '%s', '%s', '%s', '%d', '%s')",
		d.id, d.parent_id, d.name, d.path, d.size, d.mod_time)
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
