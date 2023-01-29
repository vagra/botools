package app

import (
	"gopkg.in/ini.v1"
)

func ReadConfig() bool {

	println("读取 " + CONFIG_INI)

	cfg, err := ini.Load(CONFIG_INI)
	Check(err, "读取 "+CONFIG_INI+" 失败")

	g_disks = cfg.Section("disks").KeysHash()
	if len(g_disks) <= 0 {
		println("请在 " + CONFIG_INI + " 中设置 disks 列表")
		println("格式：disk-name = disk:/path")
		return false
	}

	for name, path := range g_disks {
		println(name + " = " + path)
	}

	return true
}
