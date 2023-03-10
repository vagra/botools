package app

import (
	"fmt"
	"log"
)

func MoveErrors() error {
	println("start: move error dirs and files to special folder")

	println()
	InitLog(MOVE_ERRORS_LOG)

	println()
	ReadConfig()

	println()
	ReadDotSQL()

	println()
	ReadErrors()

	println()
	CheckErrorsRootExists()

	println()
	CheckDBsDirExists()

	println()
	CheckAllDBExists()

	println()
	CheckAllDBInited()

	println()
	CheckAllDBHasData()

	println()
	GetHasErrorDBs()

	CheckTaskHasDBs()

	println()
	ConfirmMoveErrors()

	println()
	STMoveErrors()

	println()
	println("move error dirs and files done!")
	return nil
}

func STMoveErrors() {
	fmt.Printf("复制异常文件和文件夹到 %s\n", g_roots.errors_root)

	for name := range g_dbs {
		MoveDiskErrors(name)
	}

	println()
	println("把数据库中异常 dirs 和 files 的 error 设为 1")

	for name := range g_dbs {
		MoveDBErrors(name)
	}
}

func MoveDiskErrors(disk_name string) {
	disk_path := g_disks[disk_name]

	fmt.Printf("%s worker: start copy from %s to %s\n",
		disk_name, disk_path, g_roots.errors_root)

	for _, item := range g_errors[disk_name] {
		if item.error_type == NODIR {
			continue
		}

		fmt.Printf("src  %s\n", item.RealPath())

		if !PassMakeParentDirs(item.DestPath()) {
			continue
		}

		if !PassFileExists(item.RealPath()) {
			continue
		}

		PassCopyFile(item.RealPath(), item.DestPath())
	}

	fmt.Printf("%s worker: stop.\n", disk_name)
}

func MoveDBErrors(disk_name string) {
	db_path := GetDBPath(disk_name)

	fmt.Printf("%s worker: start update %s\n", disk_name, db_path)

	db := g_dbs[disk_name]

	for _, item := range g_errors[disk_name] {
		switch item.error_type {
		case NODIR:
			MoveDBErrorDir(db_path, db, item)
		case NOFILE:
			MoveDBErrorFile(db_path, db, item)
		default:
		}
	}

	fmt.Printf("%s worker: stop.\n", disk_name)
}

func MoveDBErrorDir(db_path string, db *DB, item *ErrorItem) {
	path := item.RealPath()
	id := db.QueryDirIDFromPath(path)

	if len(id) <= 0 {
		fmt.Printf("%s no dir %s\n", db_path, path)
		log.Printf("%s no dir %s\n", db_path, path)
		return
	}

	fmt.Printf("%s   dir\n", id)

	db.ModDirError(id, 1)
	db.ModDirFilesError(id, 1)
}

func MoveDBErrorFile(db_path string, db *DB, item *ErrorItem) {
	path := item.RealPath()
	id := db.QueryFileIDFromPath(path)

	if len(id) <= 0 {
		fmt.Printf("%s no file %s\n", db_path, path)
		log.Printf("%s no file %s\n", db_path, path)
		return
	}

	fmt.Printf("%s  file\n", id)

	db.ModFileError(id, 1)
}

func ConfirmMoveErrors() {
	fmt.Printf("本程序用于把 %s 中的异常文件和文件夹复制到 errors-root 目录\n", ERRORS_TXT)
	println("请确认已做好如下准备工作：")
	fmt.Printf("1. %s 中有异常文件或文件夹的 disks 都存在，包括物理路径和数据库\n", CONFIG_INI)
	fmt.Printf("2. %s 中配置好 disks-roots 和 errors-roots\n", CONFIG_INI)
	fmt.Printf("3. errprs-root 目录若不存在，会自动创建\n")
	fmt.Printf("4. disks-root 目录不会被访问，但要用于路径换算\n")
	fmt.Printf("5. disks-root 必须与 %s 中的 disks 根目录一致，但使用正斜杠 /\n", ERRORS_TXT)
	fmt.Printf("6. %s 中的所有路径都要使用正斜杠 /\n", CONFIG_INI)
	println("您确定要执行这个操作吗？请输入 yes 或 no ：")

	CheckConfirm()
}
