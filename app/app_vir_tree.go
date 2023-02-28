package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

func VirTree() error {
	println("start: gen virtual tree")

	println()
	InitLog(VIR_TREE_LOG)

	println()
	ReadConfig()

	println()
	CheckVirDirExists()

	println()
	CheckAllRealDiskExists()

	println()
	CheckAllVirDiskExists()

	println()
	GetEmptyVirDisks()

	CheckTaskHasVirDisks()

	println()
	MTVerTree()

	println()
	println("gen virtual tree done!")
	return nil
}

func MTVerTree() {
	println("每个 disk 启动一个线程，先获取目录树，然后创建虚拟目录树")

	var wg sync.WaitGroup

	for name := range g_vdisks {
		wg.Add(1)
		go VirTreeWorker(&wg, name)
	}

	wg.Wait()
}

func VirTreeWorker(wg *sync.WaitGroup, disk_name string) {
	defer wg.Done()

	disk_path := g_disks[disk_name]
	vdisk_path := g_vdisks[disk_name]

	fmt.Printf("%s worker: start scan %s\n", disk_name, disk_path)

	VirDir(vdisk_path, disk_path)

	fmt.Printf("%s worker: stop.\n", disk_name)
}

func VirDir(vir_path string, real_path string) {

	if !DirExists(real_path) {
		log.Printf("real path not exists: %s\n", real_path)
		return
	}

	items, _ := ioutil.ReadDir(real_path)
	for _, item := range items {
		item_real_path := real_path + "/" + item.Name()

		if IsHidden(item_real_path) {
			continue
		}

		item_vir_path := vir_path + "/" + item.Name()

		if item.IsDir() {
			MakeVirDir(item_vir_path)
			VirDir(item_vir_path, item_real_path)
		} else {
			MakeSymlink(item_vir_path, item_real_path)
		}
	}
}

func MakeVirDir(vir_path string) {
	if !DirExists(vir_path) {
		err := os.Mkdir(vir_path, os.ModePerm)
		if err != nil {
			log.Printf("创建虚拟目录 %s 时出错\n", vir_path)
		}
	}
}

func MakeSymlink(vir_path string, real_path string) {
	err := os.Symlink(real_path, vir_path)
	if err != nil {
		log.Printf("创建符号链接失败：%s -> %s\n", vir_path, real_path)
	}
}
