package app

import (
	"database/sql"

	"github.com/qustavo/dotsql"
)

const STEP_COUNT int = 5
const TIME_FORMAT string = "2006-01-02 15:04:05"
const INSERT_COUNT int = 1000

const CONFIG_INI = "config.ini"

const DB_DIR string = "dbs"
const DB_EXT string = ".db"

const VIR_DIR string = "vir"

const DISKS string = "disks"
const FILES string = "files"
const DIRS string = "dirs"
const INFOS string = "infos"

const DOT_SQL string = "dot.sql"

const SQL_CHECK_TABLE string = "check-table-exists"

const SQL_CREATE_DIRS string = "create-dirs-table"
const SQL_CREATE_FILES string = "create-files-table"
const SQL_CREATE_INFOS string = "create-infos-table"

const SQL_ADD_DIR string = "add-dir"
const SQL_ADD_FILE string = "add-file"
const SQL_ADD_INFO string = "add-info"

const SQL_ADD_DIRS string = "add-dirs"
const SQL_ADD_FILES string = "add-files"

const SQL_COUNT_DIRS string = "get-dirs-count"
const SQL_COUNT_FILES string = "get-files-count"

const SQL_GET_ALL_DIRS string = "get-all-dirs"
const SQL_GET_FILES_NO_SHA1 string = "get-files-no-sha1"
const SQL_MOD_FILE_SHA1 string = "mod-file-sha1"
const SQL_MOD_FILE_STATUS string = "mod-file-status"

const SQL_TRIM_DIR_IDS string = "trim-dir-ids"
const SQL_TRIM_FILE_IDS string = "trim-file-ids"

const SQL_GET_ROOT_DIR string = "get-root-dir"
const SQL_MOD_ROOT_DIR string = "mod-root-dir"

const SQL_REPLACE_DIR_PATHS string = "replace-dir-paths"
const SQL_REPLACE_FILE_PATHS string = "replace-file-paths"

const SQL_GET_VERSION string = "get-db-version"
const SQL_MOD_VERSION string = "mod-db-version"

const GET_TREE_LOG string = "get_tree.log"
const GEN_LINK_LOG string = "gen_link.log"
const CHECKSUM_LOG string = "checksum.log"
const REAL2DB_LOG string = "real2db.log"

const MIGRATE string = "migrate-v"

const MAX_CHAN int = 100

var g_db_tables []string = []string{DIRS, FILES, INFOS}
var g_create_sqls []string = []string{SQL_CREATE_DIRS, SQL_CREATE_FILES, SQL_CREATE_INFOS}

var g_threads int
var g_disks map[string]string
var g_dbs map[string]*sql.DB
var g_dot *dotsql.DotSql

var g_map_dirs map[string]map[string]*Dir
var g_map_files map[string]map[string]*File

var g_dirs_counter map[string]*int64
var g_files_counter map[string]*int64

var g_latest int

const WELCOME string = `
BOTOOLS - bot.sanxuezang.com toolchain

请输入数字并回车来启动对应的子程序：
1)    init_db: 初始化数据库
      若数据库文件不存在就新建；若已存在就跳过。
2)    get_tree: 获取目录树
      获取指定根目录下的文件夹、文件的路径和元数据，保存到数据库；若数据库有数据就跳过。
3)    get_size: 获取文件夹大小
      基于现有数据库，获取每一个文件夹的大小。
4)    checksum: 获取文件校验和
      基于现有数据库，获取每一个文件的 SHA1 校验和。
5)    vir_tree: 生成虚拟目录树
      不生成数据库，而是用软链接的方式生成虚拟的目录树。
6)    sync_real2db: 从物理目录同步数据库
      检查物理目录的文件夹和文件，如果不存在了，在数据库中把它们的 status 标记为 1。
101)  trim_ids: 截短 ID [已禁用]
      一次性临时维护功能，数据库中的 dirs 和 files id 16 位太长，截到 8 位
102)  mod_path: 修改路径
      临时维护功能，把数据库中的 dirs 和 files 的 path 根路径替换为新的 disk 路径
103)  migrate_db: 升级数据库
      [2023-02-23 v2] 在 dirs 表添加新字段 status 用于标记文件夹状态 0存在 1不存在 2重复
0)    exit: 退出程序

请输入数字并回车来启动对应的子程序：`
