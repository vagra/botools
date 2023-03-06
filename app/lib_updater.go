package app

import (
	"botools/selfupdate"
	"fmt"
	"strings"
)

func InitUpdater() {
	g_updater = &selfupdate.Updater{
		CurrentVersion: VERSION,
		ApiURL:         URL_BASE,
		BinURL:         URL_BASE,
	}
}

func CheckVersion() {
	fmt.Printf("本地版本：%s\n", VERSION)

	remote_version := CheckUpdateAvailable()

	fmt.Printf("检查到新版本：%s\n", remote_version)
}

func CheckUpdateAvailable() string {
	version, err := g_updater.UpdateAvailable()
	Check(err, "can't get remote version")

	if len(version) <= 0 {
		println("botools 已是最新，无需更新")
		WaitExit(1)
	}

	return strings.TrimSpace(version)
}

func SelfUpdateFile(file_name string) {
	up := selfupdate.NewUpdate()
	up.TargetPath = file_name
	url := fmt.Sprintf("%s%s", URL_BASE, file_name)

	err, errRecover := up.FromUrl(url)
	Check(err, "update %s failed", file_name)
	Check(errRecover, "remote file may be missing, recover %s", file_name)

	fmt.Printf("updated: %s\n", file_name)
}

func SelfUpateExe() {
	err := g_updater.BackgroundRun()
	Check(err, "update %s failed", APP_EXE)
}
