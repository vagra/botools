package app

import (
	"fmt"
	"log"
)

func CheckMirrorsRootExists() {
	fmt.Printf("检查镜像目录 %s 是否存在\n", g_roots.mirrors_root)

	if !DirExists(g_roots.errors_root) {
		fmt.Printf("镜像目录 %s 不存在，请检查 %s\n", g_roots.mirrors_root, CONFIG_INI)
		WaitExit(1)
	}
}

func DedupDB(disk_name string) {
	db := g_dbs[disk_name]

	count := 0

	id := ""

	for {
		file, ok := db.GetNextNoDupFile(id)
		if !ok {
			break
		}

		if len(file.dup_id) > 4 {
			continue
		}

		id = file.id

		key := file.Sha1SizeKey()

		dup_id, ok := g_uniques[key]
		if !ok {
			g_uniques[key] = id

			db.ModFileDupID(file.id, "UNIQUE")

			continue
		}

		file.dup_id = dup_id
		db.ModFileDupID(file.id, file.dup_id)
		log.Printf("%s\t==\t%s\n", file.id, file.dup_id)

		count++
	}

	fmt.Printf("%s: %d dups found\n", disk_name, count)
}

func DedupMirror(disk_name string) {
	db := g_dbs[disk_name]

	total := 0
	not_exists := 0
	deleted := 0

	var id string = ""
	var path string = ""

	for {
		file, ok := db.GetNextDupFile(id)
		if !ok {
			break
		}

		id = file.id

		path, ok = file.MirrorPath()
		if !ok {
			log.Printf("file %s path in db not start with disks-root %s\n", id, g_roots.disks_root)
		}

		total++

		if !FileExists(path) {
			not_exists++
			continue
		}

		ok = RemoveFile(path)
		if !ok {
			continue
		}

		deleted++
	}

	fmt.Printf("%s: %8d dups%8d not exists%8d deleted\n",
		disk_name, total, not_exists, deleted)
}
