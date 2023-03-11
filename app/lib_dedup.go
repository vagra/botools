package app

import (
	"fmt"
	"log"
)

func InitDupMap() {
	g_dup_files = make(map[string]string)
}

func DBDedup(disk_name string) {
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
