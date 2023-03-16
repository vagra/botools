package app

import (
	"fmt"
	"os"
)

func DB2Vir(disk_name string) {

}

func CheckVirsRootExists() {
	fmt.Printf("检查虚拟文件夹 %s 是否存在\n", g_roots.virs_root)

	if !DirExists(g_roots.virs_root) {
		fmt.Printf("虚拟文件夹 %s 不存在，新建...\n", g_roots.virs_root)

		err := os.Mkdir(g_roots.virs_root, os.ModePerm)
		Check(err, "创建虚拟文件夹 %s 时出错", g_roots.virs_root)
	}
}
