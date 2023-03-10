package app

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func CheckAllRealDiskExists() {
	println("检查是否所有 real disks 路径都存在")

	if disks, yes := AllRealDiskExists(); !yes {
		println("检查到如下 real disks 的路径不存在")
		for name, path := range disks {
			fmt.Printf("%s = %s\n", name, path)
		}
		fmt.Printf("请检查 %s\n", CONFIG_INI)
		WaitExit(1)
	}
}

func AllRealDiskExists() (map[string]string, bool) {
	var disks map[string]string = make(map[string]string)

	for disk_name, disk_path := range g_disks {
		if !DirExists(disk_path) {
			disks[disk_name] = disk_path
		}
	}

	return disks, len(disks) <= 0
}

func ReadRealTree(disk_name string) {

	db := g_dbs[disk_name]

	db.ModDirsStatus(1)
	db.ModFilesStatus(1)

	root_dir := db.GetRootDir()

	ReadRealDir(db, disk_name, root_dir)
}

func ReadRealDir(db *DB, disk_name string, dir *Dir) {

	if !DirExists(dir.path) {
		log.Printf("real dir not exists: %s\n", dir.path)
		return
	}

	dir_id, ok := db.QueryDirIDFromPath(dir.path)
	if ok {
		db.ModDirStatus(dir_id, 0)
	} else {
		dir.id = GenDirUID(disk_name)
		db.AddDir(dir)
		log.Printf("new dir: %s  %s\n", dir.id, dir.path)
	}

	items, _ := os.ReadDir(dir.path)
	for _, item := range items {
		item_path := dir.path + "/" + item.Name()
		item_path = strings.Replace(item_path, "//", "/", -1)

		if IsHidden(item_path) {
			continue
		}

		if item.IsDir() {

			var sub Dir
			sub.parent_id = dir.id
			sub.name = item.Name()
			sub.path = item_path
			info, _ := item.Info()
			sub.mod_time = info.ModTime().Format(TIME_FORMAT)

			ReadRealDir(db, disk_name, &sub)

		} else {

			var file File

			file.parent_id = dir.id
			file.name = item.Name()
			file.path = item_path
			info, _ := item.Info()
			file.size = info.Size()
			file.mod_time = info.ModTime().Format(TIME_FORMAT)

			ReadRealFile(db, disk_name, &file)

		}
	}
}

func ReadRealFile(db *DB, disk_name string, file *File) {

	if !FileExists(file.path) {
		log.Printf("real file not exists: %s\n", file.path)
		return
	}

	file_id, ok := db.QueryFileIDFromPath(file.path)
	if ok {
		db.ModFileStatus(file_id, 0)
	} else {
		file.id = GenFileUID(disk_name)
		db.AddFile(file)
		log.Printf("new file: %s  %s\n", file.id, file.path)
	}
}
