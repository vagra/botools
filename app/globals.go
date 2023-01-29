package app

import (
	"database/sql"

	"github.com/qustavo/dotsql"
)

const STEP_COUNT int = 4
const DB_COUNT int = 5
const TIME_FORMAT string = "2006-01-02 15:04:05"

const CONFIG_INI = "config.ini"

const DISKS_DB string = "disks.db"
const FILES_DB string = "files.db"
const DIRS_DB string = "dirs.db"
const FILE_METAS_DB string = "file_metas.db"
const DIR_METAS_DB string = "dir_metas.db"

const DISKS_TABLE string = "disks"
const FILES_TABLE string = "files"
const DIRS_TABLE string = "dirs"
const FILE_METAS_TABLE string = "file_metas"
const DIR_METAS_TABLE string = "dir_metas"

const DOT_SQL string = "dot.sql"

const SQL_CREATE_DISKS string = "create-disks-table"
const SQL_CREATE_FILES string = "create-files-table"
const SQL_CREATE_DIRS string = "create-dirs-table"
const SQL_CREATE_FILE_METAS string = "create-file_metas-table"
const SQL_CREATE_DIR_METAS string = "create-dir_metas-table"

const SQL_ADD_DISK string = "add-disk"
const SQL_ADD_FILE string = "add-file"
const SQL_ADD_DIR string = "add-dir"
const SQL_ADD_FILE_META string = "add-file_meta"
const SQL_ADD_DIR_META string = "add-dir_meta"

const SQL_GET_ALL_DISKS string = "get-all-disks"

const SQL_GET_DISKS_COUNT string = "get-disks-count"
const SQL_GET_FILES_COUNT string = "get-files-count"
const SQL_GET_DIRS_COUNT string = "get-dirs-count"
const SQL_GET_FILE_METAS_COUNT string = "get-file_metas-count"
const SQL_GET_DIR_METAS_COUNT string = "get-dir_metas-count"

const GET_TREE_LOG string = "get_tree.log"
const GET_META_LOG string = "get_meta.log"

var g_db_names []string = []string{DISKS_DB, FILES_DB, DIRS_DB, FILE_METAS_DB, DIR_METAS_DB}
var g_db_tables []string = []string{DISKS_TABLE, FILES_TABLE, DIRS_TABLE, FILE_METAS_TABLE, DIR_METAS_TABLE}
var g_create_sqls []string = []string{SQL_CREATE_DISKS, SQL_CREATE_FILES, SQL_CREATE_DIRS, SQL_CREATE_FILE_METAS, SQL_CREATE_DIR_METAS}
var g_count_sqls []string = []string{SQL_GET_DISKS_COUNT, SQL_GET_FILES_COUNT, SQL_GET_DIRS_COUNT, SQL_GET_FILE_METAS_COUNT, SQL_GET_DIR_METAS_COUNT}

var g_disks_db *sql.DB
var g_files_db *sql.DB
var g_dirs_db *sql.DB
var g_file_metas_db *sql.DB
var g_dir_metas_db *sql.DB

var g_dbs [DB_COUNT]*sql.DB
var g_dot *dotsql.DotSql
var g_disks map[string]string

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
