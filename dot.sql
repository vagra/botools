-- name: create-dirs-table
BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS "dirs" (
	"id" VARCHAR(64) NOT NULL UNIQUE,
	"name" TEXT NOT NULL DEFAULT '',
	"parent_id" VARCHAR(64) NOT NULL DEFAULT '0',
	"size" INTEGER NOT NULL DEFAULT 0,
	"mod_time" VARCHAR(32) NOT NULL DEFAULT '',
	PRIMARY KEY("id")
);

COMMIT;

-- name: create-files-table
BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS "files" (
	"id" VARCHAR(64) NOT NULL UNIQUE,
	"name" TEXT NOT NULL DEFAULT '',
	"parent_id" VARCHAR(64) NOT NULL DEFAULT '',
	"size" INTEGER NOT NULL DEFAULT 0,
	"md5" TEXT NOT NULL DEFAULT '',
	"mod_time" VARCHAR(32) NOT NULL DEFAULT '',
	PRIMARY KEY("id")
);

COMMIT;

-- name: add-dir
INSERT INTO dirs (id, name, parent_id, mod_time) VALUES(?, ?, ?);

-- name: add-file
INSERT INTO files (id, name, parent_id, size, mod_time) VALUES(?, ?, ?, ?);

-- name: get-dirs-count
SELECT count(id) FROM dirs;

-- name: get-files-count
SELECT count(id) FROM files;