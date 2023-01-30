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
		if len(name) <= 0 || len(path) <= 0 {
			println("格式：disk-name = disk:/path")
			println("等号的左右两边不能为空")
			println("请检查 " + CONFIG_INI)
			return false
		}
		if !IsValidName(name) {
			println("disk-name 只能包含字母、数字、_、-，并且必须以字母开头")
			println("请检查 " + CONFIG_INI)
			return false
		}
		println(name + " = " + path)
	}

	return true
}
