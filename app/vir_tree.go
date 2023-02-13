package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func VirTree() error {
	println("start: gen virtual links")

	file, err := os.OpenFile(GEN_LINK_LOG, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	Check(err, "打开 "+GEN_LINK_LOG+" 时出错")
	defer file.Close()

	log.SetOutput(file)

	println()
	CheckConfig()
	CheckVirDirExist()
	CheckNoDiskDirExist()

	println()
	MakeDiskDirs()

	for name, path := range g_disks {
		GenVirTree(name, path)
		println()
	}

	println("gen virtual links done!")
	return nil
}

func CheckNoDiskDirExist() {
	if AnyDiskDirExist() {
		println()
		fmt.Printf("执行本程序会删除 %s 下已有的 disks 虚拟目录，再重新生成\n", VIR_DIR)
		fmt.Printf("如果您不想删除某个 disk 的虚拟目录，可以在 %s 中把这个 disk 注释起来\n", CONFIG_INI)
		println("您确定要删除现有的 disks 虚拟目录？请输入 yes 或 no ：")

		if Confirm() {
			println()
			DeleteDiskDirs()
		} else {
			WaitExit(1)
		}
	}
}

func GenVirTree(disk_name string, disk_path string) {
	fmt.Printf("遍历 %s : %s\n", disk_name, disk_path)

	vir_path := VIR_DIR + "/" + disk_name

	VirDir(vir_path, disk_path)
}

func VirDir(vir_path string, real_path string) {

	if !DirExist(real_path) {
		log.Printf("real path not exist: %s\n", real_path)
		return
	}

	items, _ := ioutil.ReadDir(real_path)
	for _, item := range items {
		item_real_path := real_path + "/" + item.Name()

		if IsHidden(item_real_path) {
			continue
		}

		item_vir_path := vir_path + "/" + item.Name()

		if item.IsDir() {
			MakeVirDir(item_vir_path)
			VirDir(item_vir_path, item_real_path)
		} else {
			MakeSymlink(item_vir_path, item_real_path)
		}
	}
}

func MakeVirDir(vir_path string) {
	if !DirExist(vir_path) {
		err := os.Mkdir(vir_path, os.ModePerm)
		if err != nil {
			log.Printf("创建虚拟目录 %s 时出错\n", vir_path)
		}
	}
}

func MakeSymlink(vir_path string, real_path string) {
	err := os.Symlink(real_path, vir_path)
	if err != nil {
		log.Printf("创建符号链接失败：%s -> %s\n", vir_path, real_path)
	}
}

func CheckVirDirExist() {
	if !DirExist(VIR_DIR) {
		err := os.Mkdir(VIR_DIR, os.ModePerm)
		Check(err, "创建虚拟目录 "+VIR_DIR+" 时出错")
	}
}

func CheckDiskPaths() bool {
	for name, path := range g_disks {
		if !DirExist(path) {
			fmt.Printf("%s 的路径 %s 不存在\n", name, path)
			fmt.Printf("请检查 %s\n", CONFIG_INI)
			return false
		}
	}
	return true
}

func AnyDiskDirExist() bool {
	for name := range g_disks {
		dir_name := VIR_DIR + "/" + name
		if DirExist(dir_name) {
			return true
		}
	}

	return false
}

func DeleteDiskDirs() {
	for name := range g_disks {

		dir_name := VIR_DIR + "/" + name

		fmt.Printf("删除虚拟目录 %s", dir_name)

		if !DirExist(dir_name) {
			println(" (不存在)")
			continue
		}
		println()

		err := os.RemoveAll(dir_name)
		Check(err, "删除 "+dir_name+" 目录时出错")
	}
}

func MakeDiskDirs() {
	for name := range g_disks {
		dir_name := VIR_DIR + "/" + name

		err := os.Mkdir(dir_name, os.ModePerm)
		Check(err, "创建 "+name+" 的根目录 "+dir_name+" 时出错")
	}
}
