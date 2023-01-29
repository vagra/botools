package app

type Disk struct {
	id     int64
	name   string
	path   string
	dir_id int64
	size   int64
}

type File struct {
	id        int64
	name      string
	parent_id int64
}

type FileMeta struct {
	id       int64
	file_id  int64
	size     int64
	md5      string
	mod_time string
}

type Dir struct {
	id        int64
	name      string
	parent_id int64
}

type DirMeta struct {
	id       int64
	dir_id   int64
	size     int64
	mod_time string
}
