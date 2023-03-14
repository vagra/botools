# botools
[bot.sanxuezang.com](https://bot.sanxuezang.com) toolchain

used to get the dirs/files tree on the disk, including meta, sha1, and record to the sqlite database, then deduplications, make and sync virtual links for dir and files, etc.

![chart](https://github.com/vagra/botools/blob/18ee8d109d9c337915d32cfb41e8d20ad73d1d05/assets/charts.png)

# requirement
2023-1-29

as of now, the latest stable of `go` is `1.19.5`, but this version will encounter a cgi compilation error when build this project, which is caused by `go-sqlite3`.

the solution is to use `go 1.20` or later.

# build
if your `go` is `1.20` or later, just compile and run it:
```
$ go mod tidy
$ go build
$ go run botools
```

# using
copy `example.ini` as `config.ini`, edit this `ini` file, then start `botools.exe`.
