package app

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/mattn/go-sqlite3"
)

func LoadEmptyDBs2Mem() {
	println("加载所有 dirs 和 files 表都表为空的数据库到内存")

	g_dbs = make(map[string]*DB)

	for disk_name := range g_disks {

		db_path := GetDBPath(disk_name)

		CheckDBExists(disk_name)

		db := DBOpen(db_path)

		CheckDBInited(db, db_path)

		_, yes := db.NoData()

		db.Close()
		db = nil

		if !yes {
			continue
		}

		old_path := GetOldDBPath(disk_name)
		CheckBackupDB(db_path, old_path)

		old := DBOpen(old_path)
		mem_path := GetMemDBPath(disk_name)
		mem := DBOpen(mem_path)
		CheckBakeDB(old, mem)
		println(mem_path)

		old.Close()
		old = nil

		g_dbs[disk_name] = mem
	}
}

func OnlyReadHasDataDBs2Mem() {
	println("加载所有 dirs 和 files 表都有数据的数据库到内存")

	g_dbs = make(map[string]*DB)

	for disk_name := range g_disks {

		db_path := GetDBPath(disk_name)

		CheckDBExists(disk_name)

		db := DBOpen(db_path)

		CheckDBInited(db, db_path)

		CheckDBHasData(db, db_path)

		mem_path := GetMemDBPath(disk_name)
		mem := DBOpen(mem_path)
		CheckBakeDB(db, mem)
		println(mem_path)

		db.Close()
		db = nil

		g_dbs[disk_name] = mem
	}
}

func LoadHasDataDBs2Mem() {
	println("加载所有 dirs 和 files 表都有数据的数据库到内存")

	g_dbs = make(map[string]*DB)

	for disk_name := range g_disks {

		db_path := GetDBPath(disk_name)

		CheckDBExists(disk_name)

		db := DBOpen(db_path)

		CheckDBInited(db, db_path)

		CheckDBHasData(db, db_path)

		db.Close()
		db = nil

		old_path := GetOldDBPath(disk_name)
		CheckBackupDB(db_path, old_path)

		old := DBOpen(old_path)
		mem_path := GetMemDBPath(disk_name)
		mem := DBOpen(mem_path)
		CheckBakeDB(old, mem)
		println(mem_path)

		old.Close()
		old = nil

		g_dbs[disk_name] = mem
	}
}

func BakeMemDBs() {
	println("持久化内存数据库到 dbs")

	for disk_name := range g_dbs {
		BakeMemDB(disk_name)
	}
}

func BakeMemDB(disk_name string) {

	db_path := GetDBPath(disk_name)

	fmt.Printf("bake memory db to %s\n", db_path)

	mem := g_dbs[disk_name]
	db := DBOpen(db_path)

	CheckBakeDB(mem, db)

	mem.Close()
	mem = nil

	db.Close()
	db = nil
}

func CheckBackupAllDBs() {

	for disk_name := range g_disks {
		src_path := GetDBPath(disk_name)
		dst_path := GetOldDBPath(disk_name)

		CheckBackupDB(src_path, dst_path)
	}
}

func CheckBackupDB(src_path string, dst_path string) {
	err := os.Rename(src_path, dst_path)
	Check(err, "backup db %s failed.", src_path)
}

func CheckBakeDB(src *DB, dst *DB) {
	err := BakeDB(src, dst)
	Check(err, "bake db failed.")
}

func BakeDB(src *DB, dst *DB) error {
	dst_conn, err := (*sql.DB)(dst).Conn(context.Background())
	if err != nil {
		return err
	}

	src_conn, err := (*sql.DB)(src).Conn(context.Background())
	if err != nil {
		return err
	}

	return dst_conn.Raw(func(dst_conn interface{}) error {
		return src_conn.Raw(func(src_conn interface{}) error {
			dst_lite_conn, ok := dst_conn.(*sqlite3.SQLiteConn)
			if !ok {
				return fmt.Errorf("can't convert destination connection to SQLiteConn")
			}

			src_lite_conn, ok := src_conn.(*sqlite3.SQLiteConn)
			if !ok {
				return fmt.Errorf("can't convert source connection to SQLiteConn")
			}

			b, err := dst_lite_conn.Backup("main", src_lite_conn, "main")
			if err != nil {
				return fmt.Errorf("error initializing SQLite backup: %w", err)
			}

			done, err := b.Step(-1)
			if !done {
				return fmt.Errorf("step of -1, but not done")
			}
			if err != nil {
				return fmt.Errorf("error in stepping backup: %w", err)
			}

			err = b.Finish()
			if err != nil {
				return fmt.Errorf("error finishing backup: %w", err)
			}

			return err
		})
	})
}
