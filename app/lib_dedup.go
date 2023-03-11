package app

import (
	"log"
)

func InitDupMap() {
	g_dup_files = make(map[string]string)
}

func DBDedup(disk_name string) {
	db := g_dbs[disk_name]

	id := ""

	for {
		file, ok := db.GetNextFile(id)
		if !ok {
			break
		}

		id = file.id

		key := file.Sha1SizeKey()
		// fmt.Printf("%s  %s\n", id, key)

		dup_id, ok := g_dup_files[key]
		if !ok {
			g_dup_files[key] = id
			continue
		}

		if file.dup_id == dup_id {
			continue
		}

		file.dup_id = dup_id
		db.ModFileDupID(file.id, file.dup_id)
		log.Printf("%s\t==\t%s\n", file.id, file.dup_id)
	}

}
