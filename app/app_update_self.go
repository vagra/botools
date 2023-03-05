package app

import (
	"fmt"
	"net/http"

	"github.com/minio/selfupdate"
)

func UpdateSelf() error {
	println("start: self update botools")

	println()
	CheckNewVersion()

	println()
	ConfirmUpdateSelf()

	STUpdateSelf()

	println()
	println("self update done!")
	return nil
}

func STUpdateSelf() {
	SelfUpdateFiles()
	SelfUpateExe()
}

func SelfUpdateFiles() {
	bak_sql := fmt.Sprintf("%s.old", DOT_SQL)
	PassCopyFile(DOT_SQL, bak_sql)

	bak_ini := fmt.Sprintf("%s.old", EXAMPLE_INI)
	PassCopyFile(EXAMPLE_INI, bak_ini)

	HttpDownload(DOT_SQL, URL_SQL)
	fmt.Printf("updated: %s\n", DOT_SQL)

	HttpDownload(EXAMPLE_INI, URL_INI)
	fmt.Printf("updated: %s\n", EXAMPLE_INI)
}

func SelfUpateExe() {
	resp, err := http.Get(URL_EXE)
	Check(err, "error when download file from %s", URL_EXE)
	defer resp.Body.Close()

	err = selfupdate.Apply(resp.Body, selfupdate.Options{})
	Check(err, "error when update botools.exe")

	println("updated: botools.exe")
}

func ConfirmUpdateSelf() {
	println("更新将会覆盖 botools.exe、dot.sql、config.ini.example")
	println("您确定要执行这个操作吗？请输入 yes 或 no ：")

	CheckConfirm()
}
