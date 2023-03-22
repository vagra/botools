package app

import (
	"fmt"

	"golang.org/x/exp/maps"
)

func ReadUniqueMap() {

	g_uniques = make(map[string]string)

	for _, db := range g_dbs {
		unique_files := db.GetUniqueFiles()

		for _, file := range unique_files {
			key := file.Sha1SizeKey()

			_, ok := g_uniques[key]
			if !ok {
				g_uniques[key] = file.id
			}
		}
	}

	fmt.Printf("total %d unique files in all dbs\n", len(g_unique_files))
}

func ReadUniqueOrErrorFiles() {

	g_unique_files = make(map[string]*File)

	for _, db := range g_dbs {
		unique_files := db.GetUniqueOrErrorFiles()

		maps.Copy(g_unique_files, unique_files)
	}

	fmt.Printf("total %d unique files in all dbs\n", len(g_unique_files))

}
