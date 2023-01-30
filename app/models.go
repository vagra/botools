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

type Dir struct {
	id        string
	name      string
	parent_id string
	size      int64
	mod_time  string
}

type File struct {
	id        string
	name      string
	parent_id string
	size      int64
	md5       string
	mod_time  string
}

func (d Dir) Tuple() string {
	return fmt.Sprintf("('%s', '%s', '%s', '%s')", d.id, d.name, d.parent_id, d.mod_time)
}

func (d Dir) AddMarks(marks *[]string) {
	*marks = append(*marks, "(?, ?, ?, ?)")
}

func (d Dir) AddArgs(args *[]interface{}) {
	*args = append(*args, d.id)
	*args = append(*args, d.name)
	*args = append(*args, d.parent_id)
	*args = append(*args, d.mod_time)
}

func (f File) Tuple() string {
	return fmt.Sprintf("('%s', '%s', '%s', '%d', '%s')", f.id, f.name, f.parent_id, f.size, f.mod_time)
}

func (f File) AddMarks(marks *[]string) {
	*marks = append(*marks, "(?, ?, ?, ?, ?)")
}

func (f File) AddArgs(args *[]interface{}) {
	*args = append(*args, f.id)
	*args = append(*args, f.name)
	*args = append(*args, f.parent_id)
	*args = append(*args, f.size)
	*args = append(*args, f.mod_time)
}
