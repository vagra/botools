package app

import "fmt"

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

func (e ErrorItem) Tuple() string {
	return fmt.Sprintf("%s, %d, %d, %s",
		e.disk_name, e.error_type, e.error_code, e.path)
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
