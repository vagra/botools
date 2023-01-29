-- name: create-disks-table
BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS "disks" (
	"id" INTEGER NOT NULL UNIQUE,
	"name" TEXT NOT NULL DEFAULT '',
	"path" TEXT NOT NULL DEFAULT '',
	"dir_id" INTEGER NOT NULL DEFAULT 0,
	"size" INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY("id" AUTOINCREMENT)
);

COMMIT;

-- name: create-dirs-table
BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS "dirs" (
	"id" INTEGER NOT NULL UNIQUE,
	"name" TEXT NOT NULL DEFAULT '',
	"parent_id" INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY("id" AUTOINCREMENT)
);

COMMIT;

-- name: create-files-table
BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS "files" (
	"id" INTEGER NOT NULL UNIQUE,
	"name" TEXT NOT NULL DEFAULT '',
	"parent_id" INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY("id" AUTOINCREMENT)
);

COMMIT;

-- name: create-dir_metas-table
BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS "dir_metas" (
	"id" INTEGER NOT NULL UNIQUE,
	"dir_id" INTEGER NOT NULL DEFAULT 0,
	"size" INTEGER NOT NULL DEFAULT 0,
	"mod_time" TEXT NOT NULL DEFAULT '',
	PRIMARY KEY("id" AUTOINCREMENT)
);

COMMIT;

-- name: create-file_metas-table
BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS "file_metas" (
	"id" INTEGER NOT NULL UNIQUE,
	"file_id" INTEGER NOT NULL DEFAULT 0,
	"size" INTEGER NOT NULL DEFAULT 0,
	"md5" TEXT NOT NULL DEFAULT '',
	"mod_time" TEXT NOT NULL DEFAULT '',
	PRIMARY KEY("id" AUTOINCREMENT)
);

COMMIT;

-- name: add-disk
INSERT INTO disks (name, path, dir_id) VALUES(?, ?, ?);

-- name: add-dir
INSERT INTO dirs (name, parent_id) VALUES(?, ?);

-- name: add-file
INSERT INTO files (name, parent_id) VALUES(?, ?);

-- name: add-dir_meta
INSERT INTO dir_metas (dir_id, size, mod_time) VALUES(?, ?, ?);

-- name: add-file_meta
INSERT INTO file_metas (file_id, size, mod_time) VALUES(?, ?, ?);

-- name: get-all-disks
SELECT id, name, path, dir_id FROM disks;

-- name: get-disks-count
select count(id) from disks;

-- name: get-dirs-count
select count(id) from dirs;

-- name: get-files-count
select count(id) from files;

-- name: get-dir_metas-count
select count(id) from dir_metas;

-- name: get-file_metas-count
select count(id) from file_metas;