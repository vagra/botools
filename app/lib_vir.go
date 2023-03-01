package app

import (
	"fmt"
	"os"
)

func CheckVirDirExists() {
	fmt.Printf("检查虚拟目录 %s 是否存在\n", VIR_DIR)

	if !DirExists(VIR_DIR) {
		fmt.Printf("虚拟目录 %s 不存在，新建...\n", DB_DIR)

		err := os.Mkdir(VIR_DIR, os.ModePerm)
		Check(err, "创建虚拟目录 %s 时出错", VIR_DIR)
	}
}

func GetVirDiskPath(disk_name string) string {
	return fmt.Sprintf("%s/%s", VIR_DIR, disk_name)
}

func CheckTaskHasVirDisks() {
	if len(g_vdisks) > 0 {
		return
	}

	println("没有需要处理的 vir disks")
	WaitExit(0)
}

func CheckAllVirDiskExists() {
	println("检查是否所有 vir disks 目录都存在")

	for disk_name := range g_disks {
		vdisk_path := GetVirDiskPath(disk_name)
		CheckVirDiskExists(vdisk_path)
	}
}

func GetEmptyVirDisks() {
	println("获取所有空的 vir disks")

	g_vdisks = make(map[string]string)

	for disk_name := range g_disks {

		vdisk_path := GetVirDiskPath(disk_name)

		CheckVirDiskExists(vdisk_path)

		if !VirDiskEmpty(vdisk_path) {
			continue
		}

		println(vdisk_path)

		g_vdisks[disk_name] = vdisk_path
	}
}

func CheckVirDiskExists(vdisk_path string) {
	if VirDiskExists(vdisk_path) {
		return
	}

	fmt.Printf("vir disks 目录 %s 不存在，新建...\n", vdisk_path)

	MakeVirDiskDir(vdisk_path)
}

func VirDiskExists(vdisk_path string) bool {

	return DirExists(vdisk_path)
}

func VirDiskEmpty(vdisk_path string) bool {

	items, _ := os.ReadDir(vdisk_path)

	return len(items) <= 0
}

func MakeVirDiskDir(vdisk_path string) {
	err := os.Mkdir(vdisk_path, os.ModePerm)
	Check(err, "创建 vir disks 目录 %s 时出错", vdisk_path)
}
