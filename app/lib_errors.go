package app

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func ReadErrors() {
	fmt.Printf("读取 %s\n", ERRORS_TXT)

	file, err := os.Open(ERRORS_TXT)
	Check(err, "打开 %s 时出错", ERRORS_TXT)
	defer file.Close()

	g_errors = make(map[string][]*ErrorItem)

	regex := regexp.MustCompile(ERROR_REGEX)

	scanner := bufio.NewScanner(file)
	var count int = 0
	for scanner.Scan() {
		count += 1
		line := scanner.Text()
		// fmt.Printf("%d4 %s\n", count, line)

		ScanError(regex, count, line)
	}

	err = scanner.Err()
	Check(err, "读取 %s 时出错", ERRORS_TXT)

	println("如下 disks 存在异常文件或文件夹：")
	for disk_name, errors := range g_errors {
		fmt.Printf("%s: %d\n", disk_name, len(errors))
	}
}

func ScanError(regex *regexp.Regexp, count int, line string) bool {

	line = strings.TrimSpace(line)

	if len(line) <= 0 {
		return false
	}

	matches := regex.FindStringSubmatch(line)

	item, ok := MatchError(matches)

	if !ok {
		fmt.Printf("skip: %s\n", line)
		return false
	}

	g_errors[item.disk_name] = append(g_errors[item.disk_name], item)

	return true
}

func MatchError(matches []string) (*ErrorItem, bool) {
	var item ErrorItem

	if len(matches) < 6 {
		return &item, false
	}

	item.disk_name = DiskNameFromStr(matches[4])
	item.error_type = ErrorStr2Type(matches[1])
	item.error_code = Str2Num(matches[4])
	item.path = strings.Replace(matches[5], "\\", "/", -1)

	// println(item.Tuple())

	return &item, true
}

func CheckErrorsRootExists() {
	fmt.Printf("检查异常文件夹 %s 是否存在\n", g_roots.errors_root)

	if !DirExists(g_roots.errors_root) {
		fmt.Printf("异常文件夹 %s 不存在，新建...\n", g_roots.errors_root)

		err := os.Mkdir(g_roots.errors_root, os.ModePerm)
		Check(err, "创建异常文件夹 %s 时出错", g_roots.errors_root)
	}
}
