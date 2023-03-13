package app

import (
	"fmt"
	"log"
)

func InitDupMap() {
	g_dup_files = make(map[string]string)
}

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
		file, ok := db.GetNextNodupFile(id)
		if !ok {
			break
		}

		if len(file.dup_id) > 8 {
			continue
		}

		id = file.id

		key := file.Sha1SizeKey()
		// fmt.Printf("%s  %s\n", id, key)

		dup_id, ok := g_dup_files[key]
		if !ok {
			g_dup_files[key] = id
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

	count := 0

	id := ""

	for {
		file, ok := db.GetNextNodupFile(id)
		if !ok {
			break
		}

		if len(file.dup_id) > 8 {
			continue
		}

		id = file.id

		key := file.Sha1SizeKey()
		// fmt.Printf("%s  %s\n", id, key)

		dup_id, ok := g_dup_files[key]
		if !ok {
			g_dup_files[key] = id
			continue
		}

		file.dup_id = dup_id
		db.ModFileDupID(file.id, file.dup_id)
		log.Printf("%s\t==\t%s\n", file.id, file.dup_id)

		count++
	}

	fmt.Printf("%s: %d dups found\n", disk_name, count)
}
