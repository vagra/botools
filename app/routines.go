package app

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

func Checker(wg *sync.WaitGroup, ctx context.Context, i int, ci <-chan *File, co chan<- *File) {

	wg.Add(1)
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("checker %d : stop.\n", i)
			return

		case file := <-ci:

			if file == nil {
				fmt.Printf("\ninChan -> checker %d -> outChan: no more files.\n", i)

				co <- nil

				continue
			}

			sha1, code := SHA1(file.path)
			file.sha1 = sha1
			file.status = code

			co <- file

		default:
			time.Sleep(time.Millisecond)
		}
	}
}

func Writer(wg *sync.WaitGroup, ctx context.Context, db *sql.DB, co <-chan *File, ce chan<- bool) {
	wg.Add(1)
	defer wg.Done()

	total := len(g_map_files)

	divisor := int(total / 50)

	println("__________________________________________________")

	count := 0
	for {
		select {
		case <-ctx.Done():
			println("writer: stop.")
			return

		case file := <-co:

			if file == nil {
				println("outChan -> writer -> endChan: no more files.\n")
				ce <- true

				continue
			}

			// fmt.Printf("%d/%d\twriter <- outChain: %d  %x  %s\n", count, total, file.status, file.sha1, file.path)

			UpdateFile(db, file)

			count++

			if count%divisor == 0 {
				print(".")
			}
			if count >= total {
				println()
			}

		default:
			// println("writer: waiting for signal...")
			time.Sleep(time.Millisecond)
		}
	}
}

func SHA1(path string) (string, int8) {

	file, err := os.Open(path)
	if err != nil {
		return "", 1
	}

	sha1h := sha1.New()
	io.Copy(sha1h, file)
	sum := hex.EncodeToString(sha1h.Sum(nil))

	return sum, 0
}

func UpdateFile(db *sql.DB, file *File) {
	if file.status == 0 {
		UpdateFileSha1(db, file)
	} else {
		UpdateFileStatus(db, file)
	}
}

func UpdateFileSha1(db *sql.DB, file *File) {

	_, err := g_dot.Exec(db, SQL_MOD_FILE_SHA1, file.sha1, file.id)
	if err != nil {
		log.Printf("db update file sha1 error: %s\n", err.Error())
	}

}

func UpdateFileStatus(db *sql.DB, file *File) {

	_, err := g_dot.Exec(db, SQL_MOD_FILE_STATUS, file.status, file.id)
	if err != nil {
		log.Printf("db update file status error: %s\n" + err.Error())
	}
}
