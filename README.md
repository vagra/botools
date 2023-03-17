# botools
[bot.sanxuezang.com](https://bot.sanxuezang.com) toolchain

used to get the dirs/files tree on the disk, including meta, sha1, and record to the sqlite database, then deduplications, make and sync virtual links for dir and files, etc.

![charts](https://github.com/vagra/botools/blob/8439b0413012a7d7a9d748163e57d3c0d9db2324/assets/charts.png)

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
