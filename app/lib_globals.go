package app

import (
	"botools/selfupdate"

	"github.com/qustavo/dotsql"
)

const VERSION = "1.9.2"

const TIME_FORMAT string = "2006-01-02 15:04:05"
const INSERT_COUNT int = 1000

const EXAMPLE_INI = "example.ini"
const CONFIG_INI = "config.ini"
const ERRORS_TXT = "errors.txt"
const DOT_SQL string = "dot.sql"
const APP_EXE string = "botools.exe"

const URL_BASE = "https://kepan.org/botools/"

const DISK_PRE = "disk-"
const DIR_PRE = "d-"
const FILE_PRE = "f-"

const DB_DIR string = "dbs"
const DB_EXT string = ".db"

const VIR_DIR string = "vir"

const DISKS string = "disks"
const FILES string = "files"
const DIRS string = "dirs"
const INFOS string = "infos"

const GET_TREE_LOG string = "get_tree.log"
const VIR_TREE_LOG string = "vir_tree.log"
const CHECKSUM_LOG string = "checksum.log"
const REAL2DB_LOG string = "real2db.log"
const DEDUP_DBS_LOG string = "dedup_dbs.log"
const DEDUP_MIRRORS_LOG string = "dedup_mirrors.log"
const DB2VIR_LOG string = "db2vir.log"
const MOVE_ERRORS_LOG string = "move_errors.log"

const MIGRATE string = "migrate-v"

const MAX_CHAN int = 1000

const ERROR_REGEX = `^(.*?)\((.*?)(\d+)\) : X:\\disks\\(\d+)\\(.*?)$`
const ID_REGEX = `^(.*?)(\d{8})$`

var g_disks map[string]string
var g_dbs map[string]*DB
var g_dot *dotsql.DotSql
var g_vdisks map[string]string

var g_roots *Roots

var g_threads int

var g_errors map[string][]*ErrorItem

var g_map_dirs map[string]map[string]*Dir
var g_map_files map[string]map[string]*File

var g_dirs_counter map[string]*int64
var g_files_counter map[string]*int64

var g_dup_files map[string]string

var g_real_files map[string]*File

var g_latest int

var g_updater *selfupdate.Updater

const WELCOME string = `
BOTOOLS %s - bot.sanxuezang.com toolchain

请输入数字并回车来启动对应的子程序：
1)    init_db: 初始化数据库
      若数据库文件不存在就新建；若已存在就跳过。
2)    get_tree: 获取目录树
      获取指定根目录下的文件夹、文件的路径和元数据，保存到数据库；若数据库有数据就跳过。
4)    checksum: 获取文件校验和
      基于现有数据库，获取每一个文件的 SHA1 校验和。
6)    sync_real2db: 从物理目录同步数据库
      检查物理目录的文件夹和文件，更新数据库中的 dirs, files。
7)    dedup_dbs: 在数据库中查重
      检查数据库中的 files，将重复文件的 dup_id 设为唯一文件的 id。
8)    dedup_mirrors: 在镜像目录下查重
      根据查重后的数据库，删除镜像目录下所有的重复文件，只保留一个唯一文件。
9)    sync_db2vir: 从数据库同步到虚拟目录树
      根据数据库中的 dirs 和 files 同步虚拟目录树。

100)  update_self: 更新 botools
      自动查询远程版本，比当前版本新就下载并热更新，包括 exe、dot.sql 和 example.ini 。
102)  mod_path: 修改路径
      维护功能，把数据库中的 dirs 和 files 的 path 根路径替换为新的 disk 路径。
103)  move_errors: 复制异常文件和文件夹到指定目录
      维护功能，把名字或路径超长，或包含特殊字符的文件和文件夹复制到 errors-root。
104)  mod_disk_ids: 修改数据库中的 disk_id
      维护功能，修改了 config.ini 和 dbs 的 disk_name 后，更新数据库中的 dirs 和 files 的 disk_id 。

200)  migrate_db: 升级数据库
      [2023-03-03 v4] 为 dirs 和 files 添加 error, dup_id 以标记异常、重复，status 仅用于标记是否存在
      [2023-03-02 v3] 新建表 vdirs 和 vfiles ，用于在数据库中生成虚拟树（vdb）
      [2023-02-23 v2] 在 dirs 表添加新字段 status 用于标记文件夹状态 0存在 1不存在 2重复 3名字超长
0)    exit: 退出程序

请输入数字并回车来启动对应的子程序：
`
