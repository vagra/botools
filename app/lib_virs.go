package app

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gotopkg/mslnk/pkg/mslnk"
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

		vir_path, ok := file.VirPath()
		if !ok {
			log.Printf("file %s get vir path error\n", id)
		}

		vir_path = fmt.Sprintf("%s.lnk", vir_path)

		if !PassMakeParentDirs(vir_path) {
			continue
		}

		if !MakeVirLink(real_path, vir_path, id) {
			continue
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

func MakeVirLink(real_path string, vir_path string, id string) bool {

	h := mslnk.Header()
	x := mslnk.StringData{
		"NameString": mslnk.StringDataStruct(id),
	}
	x.Update(&h)

	err := LinkFile(real_path, vir_path, id)
	if err != nil {
		log.Printf("make vir link fail: %s -> %s\n", vir_path, real_path)
		return false
	}

	return true
}

func LinkFile(target string, name string, id string) error {
	target = strings.TrimSpace(target)
	var drive string
	if len(target) < 3 || !(target[1] == ':' && target[2] == '\\') {
		drive = "C:\\"
	} else {
		drive = target[:3]
		target = target[3:]
	}

	r := mslnk.ShellLink{
		ShellLinkHeader: mslnk.Header(),
		LinkTargetIDList: mslnk.LinkTargetIDList{
			ItemIDList: []mslnk.ItemID{
				mslnk.ItemIDCLSID(mslnk.ItemIDMagic["MY_COMPUTER"]),
				mslnk.ItemIDDrive(drive),
				mslnk.ItemIDFile(target),
			},
		},
		StringData: mslnk.StringData{
			"NameString": mslnk.StringDataStruct(id),
		},
	}

	r.LinkTargetIDList.Size()
	r.ShellLinkHeader.LinkFlags["HasLinkTargetIDList"] = true
	r.ShellLinkHeader.LinkFlags["ForceNoLinkInfo"] = true
	r.ShellLinkHeader.LinkFlags["HasName"] = true
	r.ShellLinkHeader.FileAttributes["FILE_ATTRIBUTE_NORMAL"] = true
	r.ShellLinkHeader.Update()

	return r.Save(name)
}
