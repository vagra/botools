package app

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"
)

func Checker(wg *sync.WaitGroup, ctx context.Context, disk_name string, i int, ci <-chan *File, co chan<- *File) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s: checker %d stop.\n", disk_name, i)
			return

		case file := <-ci:

			if file == nil {
				fmt.Printf("%s: inChan -> checker %d -> outChan: no more files.\n", disk_name, i)

				co <- nil

				continue
			}

			sha1, code := GetSHA1(file.path)
			file.sha1 = sha1
			file.status = code

			co <- file

		default:
			time.Sleep(time.Millisecond)
		}
	}
}

func Writer(wg *sync.WaitGroup, ctx context.Context, disk_name string, co <-chan *File, ce chan<- bool) {
	defer wg.Done()

	var db *sql.DB = g_dbs[disk_name]

	total := len(g_map_files[disk_name])
	divisor := int(total / 20)

	count := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s: writer stop.\n", disk_name)
			return

		case file := <-co:

			if file == nil {
				fmt.Printf("%s: outChan -> writer -> endChan: no more files.\n", disk_name)
				ce <- true

				continue
			}

			DBUpdateFile(db, file)

			count++

			if count%divisor == 0 {
				fmt.Printf("%s: %d%%\n", disk_name, count*100/total+1)
			}

		default:
			// time.Sleep(time.Millisecond)
		}
	}
}
