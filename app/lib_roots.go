package app

import "fmt"

type Roots struct {
	disks_root  string
	errors_root string
	dups_root   string
	virs_root   string
}

func (r Roots) Tuple() string {
	return fmt.Sprintln("[roots]") +
		fmt.Sprintf("disks-root = %s\n", r.disks_root) +
		fmt.Sprintf("errors-root = %s\n", r.errors_root) +
		fmt.Sprintf("dups-root = %s\n", r.dups_root) +
		fmt.Sprintf("virs-root = %s", r.virs_root)
}

func CheckHasRootsConfig() {
	if !HasRootsConfig() {
		fmt.Printf("请在 %s 中设置 roots 列表，格式：\n", CONFIG_INI)
		println("[roots]")
		println("disks-root  = disk:/path/  # 所有 disks 的根目录")
		println("errors-root = disk:/path/  # 用于放置异常文件或文件夹")
		println("dups-root   = disk:/path/  # 用于放置重复文件")
		println("virs-root   = disk:/path/  # 用于放置虚拟树")
		WaitExit(1)
	}
}

func HasRootsConfig() bool {
	return DiskPathValid(g_roots.disks_root) &&
		DiskPathValid(g_roots.errors_root)
}
