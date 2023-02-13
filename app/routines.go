package app

import (
	"context"
	"database/sql"
	"fmt"
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

func Writer(wg *sync.WaitGroup, ctx context.Context, disk_name string, co <-chan *File, ce chan<- bool) {
	wg.Add(1)
	defer wg.Done()

	var db *sql.DB = g_dbs[disk_name]

	total := len(g_map_files[disk_name])
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

			DBUpdateFile(db, file)

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
