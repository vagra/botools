package app

import (
	"fmt"
	"strings"
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

func (f *File) MirrorPath() (string, bool) {
	var path string = ""

	if !strings.HasPrefix(f.path, g_roots.disks_root) {
		return f.path, false
	}

	path = fmt.Sprintf("%s/%s",
		g_roots.mirrors_root, strings.TrimLeft(f.path, g_roots.disks_root))

	path = strings.ReplaceAll(path, "//", "/")

	return path, true
}

func (f *File) RealPath() (string, bool) {

	if f.error == 1 {
		return f.ErrorPath()
	}

	if len(f.dup_id) <= 0 {
		return f.path, true
	}

	real_file, ok := g_unique_files[f.dup_id]
	if !ok {
		return f.path, false
	}

	return real_file.path, true
}

func (f *File) VirPath() (string, bool) {
	var path string = ""

	if !strings.HasPrefix(f.path, g_roots.disks_root) {
		return f.path, false
	}

	path = fmt.Sprintf("%s/%s",
		g_roots.virs_root, strings.TrimLeft(f.path, g_roots.disks_root))

	path = strings.ReplaceAll(path, "//", "/")

	return path, true
}

func (f *File) ErrorPath() (string, bool) {
	var path string = ""

	if !strings.HasPrefix(f.path, g_roots.disks_root) {
		return f.path, false
	}

	path = fmt.Sprintf("%s/%s",
		g_roots.errors_root, strings.TrimLeft(f.path, g_roots.disks_root))

	path = strings.ReplaceAll(path, "//", "/")

	return path, true
}
