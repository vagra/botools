package app

import (
	"fmt"
)

type ErrorType int

const (
	OTHER  ErrorType = iota
	NODIR  ErrorType = iota
	NOFILE ErrorType = iota
)

type ErrorItem struct {
	disk_name  string
	error_type ErrorType
	error_code int
	path       string
}

func (e *ErrorItem) Tuple() string {
	return fmt.Sprintf("%s, %d, %d, %s",
		e.disk_name, e.error_type, e.error_code, e.path)
}

func (e *ErrorItem) DiskCode() string {
	return DiskCodeFromName(e.disk_name)
}

func (e *ErrorItem) ErrorRoot() string {
	root := fmt.Sprintf("%s/%s/", g_roots.disks_root, e.DiskCode())
	root = CleanPath(root)

	return root
}

func (e *ErrorItem) RealRoot() string {
	root := fmt.Sprintf("%s/", g_disks[e.disk_name])
	root = CleanPath(root)

	return root
}

func (e *ErrorItem) DestRoot() string {

	root := fmt.Sprintf("%s/%s/", g_roots.errors_root, e.DiskCode())
	root = CleanPath(root)

	return root
}

func (e *ErrorItem) RealPath() string {

	path := fmt.Sprintf("%s/%s", e.RealRoot(), e.path)
	path = CleanPath(path)

	return path
}

func (e *ErrorItem) DestPath() string {

	path := fmt.Sprintf("%s/%s", e.DestRoot(), e.path)
	path = CleanPath(path)

	return path
}

func ErrorStr2Type(str string) ErrorType {
	switch str {
	case "CreateFile":
		return NOFILE
	case "CreateDirectory":
		return NODIR
	default:
		return OTHER
	}
}
