package app

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/mattn/go-sqlite3"
)

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

func BakeDB(dst *DB, src *DB) error {
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
