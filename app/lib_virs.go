package app

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/exp/maps"
)

func DB2Vir(disk_name string) {
	db := g_dbs[disk_name]

	count := 0

	var id string = ""

	for {
		file, ok := db.GetNextExistOrErrorFile(id)
		if !ok {
			break
		}

		id = file.id

		real_path, ok := file.RealPath()
		if !ok {
			log.Printf("file %s get real path error\n", id)
		}

		if file.error == 1 {
			fmt.Printf("%s: dup %s, error %d, real %s\n", id, file.dup_id, file.error, real_path)
		}

		count++
	}

	fmt.Printf("%s: %d files\n", disk_name, count)
}

func CheckVirsRootExists() {
	fmt.Printf("检查虚拟文件夹 %s 是否存在\n", g_roots.virs_root)

	if !DirExists(g_roots.virs_root) {
		fmt.Printf("虚拟文件夹 %s 不存在，新建...\n", g_roots.virs_root)

		err := os.Mkdir(g_roots.virs_root, os.ModePerm)
		Check(err, "创建虚拟文件夹 %s 时出错", g_roots.virs_root)
	}
}

func ReadRealMap() {

	g_real_files = make(map[string]*File)

	for _, db := range g_dbs {
		real_files := db.GetRealFiles()

		maps.Copy(g_real_files, real_files)
	}

	fmt.Printf("total %d real files in all dbs\n", len(g_real_files))

}

func MakeVirLink(file *File) {

}
