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
	error     int8
	mod_time  string
}

type File struct {
	id        string
	parent_id string
	name      string
	path      string
	size      int64
	status    int8
	error     int8
	dup_id    string
	sha1      string
	mod_time  string
}

type VDir struct {
	id        string
	parent_id string
	name      string
	path      string
	status    int8
	mod_time  string
}

type VFile struct {
	id        string
	real_id   string
	parent_id string
	name      string
	path      string
	status    int8
	mod_time  string
}

// --------------------------------------------
// general methods
// --------------------------------------------

func (d *Disk) Tuple() string {
	return fmt.Sprintf("('%d', '%s', '%s', '%d')",
		d.id, d.name, d.path, d.size)
}

func (i *Info) Tuple() string {
	return fmt.Sprintf("('%d', '%d')",
		i.id, i.db_version)
}

func (d *Dir) Tuple() string {
	return fmt.Sprintf("('%s', '%s', '%s', '%s', '%d', '%d', '%d', '%s')",
		d.id, d.parent_id, d.name, d.path, d.size, d.status, d.error, d.mod_time)
}

func (f *File) Tuple() string {
	return fmt.Sprintf("('%s', '%s', '%s', '%s', '%d', '%d', '%d', '%s', '%s', '%s')",
		f.id, f.parent_id, f.name, f.path, f.size,
		f.status, f.error, f.dup_id, f.sha1, f.mod_time)
}

func (d *VDir) Tuple() string {
	return fmt.Sprintf("('%s', '%s', '%s', '%s', '%d', '%s')",
		d.id, d.parent_id, d.name, d.path, d.status, d.mod_time)
}

func (f *VFile) Tuple() string {
	return fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%d', '%s')",
		f.id, f.real_id, f.parent_id, f.name, f.path, f.status, f.mod_time)
}

// --------------------------------------------
// methods of File
// --------------------------------------------

func (f *File) Sha1SizeKey() string {
	return fmt.Sprintf("%s-%d", f.sha1, f.size)
}
