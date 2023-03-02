package app

import (
	"fmt"

	"gopkg.in/ini.v1"
)

func ReadConfig() {
	fmt.Printf("读取 %s\n", CONFIG_INI)

	cfg, err := ini.Load(CONFIG_INI)
	Check(err, "读取 %s 失败", CONFIG_INI)

	g_disks = cfg.Section("disks").KeysHash()
	if len(g_disks) <= 0 {
		fmt.Printf("请在 %s 中设置 disks 列表\n", CONFIG_INI)
		println("格式：")
		println("[disks]")
		println("disk-name = disk:/path/")
		println("disk-name = disk:/path/")
		WaitExit(1)
	}

	println("[disks]")

	for name, path := range g_disks {
		if len(name) <= 0 || len(path) <= 0 {
			println("格式：disk-name = disk:/path")
			println("等号的左右两边不能为空")
			fmt.Printf("请检查 %s\n", CONFIG_INI)
			WaitExit(1)
		}

		if !IsValidName(name) {
			println("disk-name 只能包含字母、数字、_、-，并且必须以字母开头")
			fmt.Printf("请检查 %s\n", CONFIG_INI)
			WaitExit(1)
		}

		fmt.Printf("%s = %s\n", name, path)
	}

	g_roots = &Roots{}
	g_roots.disks_root = cfg.Section("roots").Key("disks-root").String()
	g_roots.errors_root = cfg.Section("roots").Key("errors-root").String()
	g_roots.dups_root = cfg.Section("roots").Key("dups-root").String()
	g_roots.virs_root = cfg.Section("roots").Key("virs-root").String()

	println(g_roots.Tuple())

	g_threads, err = cfg.Section("checksum").Key("threads").Int()
	if err != nil {
		fmt.Printf("请在 %s 中设置 checksum 时的线程数\n", CONFIG_INI)
		println("格式：")
		println("[checksum]")
		println("threads = 10")
		WaitExit(1)
	}

	println("[threads]")
	fmt.Printf("threads = %d\n", g_threads)

}
