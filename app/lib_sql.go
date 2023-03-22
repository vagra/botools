package app

// -----------------------------------------------
// CREATE TABLES
// -----------------------------------------------

const SQL_CREATE_DIRS string = "create-dirs-table"
const SQL_CREATE_FILES string = "create-files-table"
const SQL_CREATE_VDIRS string = "create-vdirs-table"
const SQL_CREATE_VFILES string = "create-vfiles-table"
const SQL_CREATE_INFOS string = "create-infos-table"

// -----------------------------------------------
// SPECIALES
// -----------------------------------------------

const SQL_INIT_PRAGMA string = "init-pragma"
const SQL_BEGIN string = "begin"
const SQL_END string = "end"
const SQL_CHECK_TABLE string = "check-table-exists"

// -----------------------------------------------
// INSERT
// -----------------------------------------------

//// single insert
const SQL_ADD_DIR string = "add-dir"
const SQL_ADD_FILE string = "add-file"
const SQL_ADD_INFO string = "add-info"

//// batch insert
const SQL_ADD_DIRS string = "add-dirs"
const SQL_ADD_FILES string = "add-files"

// -----------------------------------------------
// QUERY
// -----------------------------------------------

//// query dirs
const SQL_GET_ROOT_DIR string = "get-root-dir"
const SQL_COUNT_DIRS string = "get-dirs-count"
const SQL_MAX_DIR_ID string = "get-max-dir-id"
const SQL_GET_ALL_DIRS string = "get-all-dirs"
const SQL_PATH_GET_DIR_ID string = "get-dir-id-from-path"
const SQL_GET_A_DIR_ID string = "get-a-dir-id"
const SQL_GET_NEXT_DIR string = "get-next-dir"

//// query files
const SQL_COUNT_FILES string = "get-files-count"
const SQL_MAX_FILE_ID string = "get-max-file-id"
const SQL_GET_ALL_FILES string = "get-all-files"
const SQL_GET_UNIQUE_FILES string = "get-unique-files"
const SQL_GET_UNIQUE_OR_ERROR_FILES string = "get-unique-or-error-files"
const SQL_GET_NO_SHA1_FILES_COUNT string = "get-no-sha1-files-count"
const SQL_GET_NO_SHA1_FILES string = "get-no-sha1-files"
const SQL_PATH_GET_FILE_ID string = "get-file-id-from-path"
const SQL_GET_A_FILE_ID string = "get-a-file-id"
const SQL_GET_NEXT_FILE string = "get-next-file"
const SQL_GET_NEXT_NODUP_FILE string = "get-next-nodup-file"
const SQL_GET_NEXT_DUP_FILE string = "get-next-dup-file"
const SQL_GET_NEXT_EXIST_OR_ERROR_FILE string = "get-next-exist-or-error-file"

//// query infos
const SQL_GET_VERSION string = "get-db-version"

// -----------------------------------------------
// UPDATE
// -----------------------------------------------

//// update dirs
const SQL_MOD_ROOT_DIR string = "mod-root-dir"
const SQL_TRIM_DIRS_ID string = "trim-dirs-id"
const SQL_MOD_DIRS_STATUS string = "mod-dirs-status"
const SQL_REPLACE_DIRS_PATH string = "replace-dirs-path"
const SQL_REPLACE_DIRS_ID string = "replace-dirs-id"
const SQL_REPLACE_DIRS_PARENT_ID string = "replace-dirs-parent-id"
const SQL_MOD_DIR_ERROR string = "mod-dir-error"
const SQL_MOD_DIR_STATUS string = "mod-dir-status"

//// update files
const SQL_TRIM_FILES_ID string = "trim-files-id"
const SQL_MOD_FILES_STATUS string = "mod-files-status"
const SQL_RESET_FILES_DUP_ID string = "reset-files-dup-id"
const SQL_REPLACE_FILES_PATH string = "replace-files-path"
const SQL_REPLACE_FILES_ID string = "replace-files-id"
const SQL_REPLACE_FILES_PARENT_ID string = "replace-files-parent-id"
const SQL_MOD_DIR_FILES_ERROR string = "mod-dir-files-error"
const SQL_MOD_FILE_SHA1 string = "mod-file-sha1"
const SQL_MOD_FILE_STATUS string = "mod-file-status"
const SQL_MOD_FILE_ERROR string = "mod-file-error"
const SQL_MOD_FILE_DUP_ID string = "mod-file-dup-id"

//// update infos
const SQL_MOD_VERSION string = "mod-db-version"
