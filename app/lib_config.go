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
	CheckHasDisksConfig()

	println("[disks]")

	for name, path := range g_disks {

		CheckDiskNameValid(name)
		CheckDiskPathValid(path)

		fmt.Printf("%s = %s\n", name, path)
	}

	g_roots = &Roots{}
	g_roots.disks_root = cfg.Section("roots").Key("disks-root").String()
	g_roots.errors_root = cfg.Section("roots").Key("errors-root").String()
	g_roots.mirrors_root = cfg.Section("roots").Key("mirrors-root").String()
	g_roots.virs_root = cfg.Section("roots").Key("virs-root").String()

	CheckHasRootsConfig()

	println(g_roots.Tuple())

	g_threads, err = cfg.Section("checksum").Key("threads").Int()
	Check(err,
		fmt.Sprintf("请在 %s 中设置 checksum 时的线程数，格式：\n[checksum]\nthreads = 10", CONFIG_INI))

	println("[threads]")
	fmt.Printf("threads = %d\n", g_threads)

}
