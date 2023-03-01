package app

import (
	"fmt"
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

func (d Disk) Tuple() string {
	return fmt.Sprintf("('%d', '%s', '%s', '%d')",
		d.id, d.name, d.path, d.size)
}

func (i Info) Tuple() string {
	return fmt.Sprintf("('%d', '%d')",
		i.id, i.db_version)
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
