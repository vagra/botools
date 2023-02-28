package app

import (
	"fmt"
	"log"
	"os"
)

func InitLog(file_path string) {
	fmt.Printf("初始化日志 %s\n", file_path)

	file, err := os.OpenFile(file_path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	Check(err, "打开日志文件 %s 时出错", file_path)

	log.SetOutput(file)
}
