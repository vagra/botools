# botools
[bot.sanxuezang.com](https://bot.sanxuezang.com) toolchain

used to get the dirs/files tree on the disk, including meta, md5, and record to the sqlite database.

# requirement
[2023-1-29]

as of now, the latest stable of `go` is `1.19.5`, but this version will encounter a cgi compilation error when build this project, which is caused by `go-sqlite3`.

the solution is to use `go 1.20rc3` or later, which can be installed with the following command:
```
$ go install golang.org/dl/go1.20rc3@latest
$ go1.20rc3 download
```
then use the following commands to build and run:
```
$ go1.20rc3 mod tidy
$ go1.20rc3 build
$ go1.20rc3 run botools
```

# build
```
$ go mod tidy
$ go build
$ go run botools
```

# using
copy config.ini.example as config.ini, edit this ini, then start botools.
