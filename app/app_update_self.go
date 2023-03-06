package app

func UpdateSelf() error {
	println("start: self update botools")

	InitUpdater()

	println()
	CheckVersion()

	println()
	ConfirmUpdateSelf()

	STUpdateSelf()

	println()
	println("self update done!")
	return nil
}

func STUpdateSelf() {
	SelfUpdateFile(DOT_SQL)
	SelfUpdateFile(EXAMPLE_INI)
	SelfUpateExe()
}

func ConfirmUpdateSelf() {
	println("更新会自动备份然后覆盖 botools.exe、dot.sql、example.ini")
	println("您确定要执行这个操作吗？请输入 yes 或 no ：")

	CheckConfirm()
}
