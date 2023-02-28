package app

import "fmt"

func CheckAllRealDiskExists() {
	println("检查是否所有的 real disks 路径都存在")

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
