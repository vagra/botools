package app

import (
	"database/sql"

	"github.com/qustavo/dotsql"
)

const STEP_COUNT int = 4
const TIME_FORMAT string = "2006-01-02 15:04:05"
const INSERT_COUNT int = 1000

const CONFIG_INI = "config.ini"

const DB_DIR string = "dbs"
const DB_EXT string = ".db"

const DISKS string = "disks"
const FILES string = "files"
const DIRS string = "dirs"

const DOT_SQL string = "dot.sql"

const SQL_CREATE_DIRS string = "create-dirs-table"
const SQL_CREATE_FILES string = "create-files-table"

const SQL_ADD_DIR string = "add-dir"
const SQL_ADD_FILE string = "add-file"

const SQL_ADD_DIRS string = "INSERT INTO dirs (id, name, parent_id, mod_time) VALUES\n"
const SQL_ADD_FILES string = "INSERT INTO files (id, name, parent_id, size, mod_time) VALUES\n"

const SQL_COUNT_DIRS string = "get-dirs-count"
const SQL_COUNT_FILES string = "get-files-count"

const GET_TREE_LOG string = "get_tree.log"

var g_db_tables []string = []string{DIRS, FILES}
var g_create_sqls []string = []string{SQL_CREATE_DIRS, SQL_CREATE_FILES}
var g_count_sqls []string = []string{SQL_COUNT_DIRS, SQL_COUNT_FILES}

var g_disks map[string]string
var g_dbs map[string]*sql.DB
var g_dot *dotsql.DotSql

var g_map_dirs map[string]Dir
var g_map_files map[string]File

var g_dirs_counter int64
var g_files_counter int64

const WELCOME string = `
BOTOOLS - bot.sanxuezang.com toolchain

请输入数字并回车来启动对应的子程序：
1)  init_db: 初始化数据库
    若数据库文件不存在就新建；若已存在，则会删除再重建，慎重。
2)  get_tree: 获取目录树
    获取指定根目录下的文件夹、文件的路径和元数据，保存到数据库。
3)  get_size: 获取文件夹大小
    基于现有数据库，获取每一个文件夹的大小。
4)  get_md5: 获取文件 MD5
    基于现有数据库，获取每一个文件的MD5。
0)  exit: 退出程序

请输入数字并回车来启动对应的子程序：`
