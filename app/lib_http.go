package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func CheckNewVersion() {
	fmt.Printf("本地版本：%s\n", VERSION)

	g_remote_version = HttpGetVersion()

	if !IsNewVersion(VERSION, g_remote_version) {
		println("botools 已是最新，无需更新")
		WaitExit(1)
	}

	fmt.Printf("检查到新版本：%s\n", g_remote_version)
}

func HttpGetVersion() string {
	resp, err := http.Get(URL_VERSION)
	Check(err, "can't visit version url")

	data, err := io.ReadAll(resp.Body)
	Check(err, "can't read version content")

	version := string(data)
	if len(version) <= 0 {
		println("received blank string")
		WaitExit(1)
	}

	return string(data)
}

func HttpDownload(path string, url string) {
	out, err := os.Create(path)
	Check(err, "error when create file %s", path)
	defer out.Close()

	resp, err := http.Get(url)
	Check(err, "error when download file from %s", url)
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	Check(err, "error when save as file %s", path)
}
