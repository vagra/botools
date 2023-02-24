-- name: create-dirs-table
BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS "dirs" (
	"id" VARCHAR(32) NOT NULL UNIQUE,
	"parent_id" VARCHAR(32) NOT NULL DEFAULT '0',
	"name" TEXT NOT NULL DEFAULT '',
	"path" TEXT NOT NULL DEFAULT '',
	"size" INTEGER NOT NULL DEFAULT 0,
	"mod_time" VARCHAR(32) NOT NULL DEFAULT '',
	PRIMARY KEY("id")
);

COMMIT;

-- name: create-files-table
BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS "files" (
	"id" VARCHAR(32) NOT NULL UNIQUE,
	"parent_id" VARCHAR(32) NOT NULL DEFAULT '',
	"name" TEXT NOT NULL DEFAULT '',
	"path" TEXT NOT NULL DEFAULT '',
	"size" INTEGER NOT NULL DEFAULT 0,
	"status" INTEGER NOT NULL DEFAULT 0,
	"sha1" VARCHAR(64) NOT NULL DEFAULT '',
	"mod_time" VARCHAR(32) NOT NULL DEFAULT '',
	PRIMARY KEY("id")
);

COMMIT;

-- name: create-infos-table
BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS "infos" (
	"id" INTEGER NOT NULL UNIQUE,
	"db_version" INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY("id" AUTOINCREMENT)
);

COMMIT;

INSERT INTO infos (db_version) VALUES(1);

-- name: add-dir
INSERT INTO dirs (id, parent_id, name, path, mod_time) VALUES(?, ?, ?, ?);

-- name: add-file
INSERT INTO files (id, parent_id, name, path, size, mod_time) VALUES(?, ?, ?, ?, ?);

-- name: add-info
INSERT INTO infos (db_version) VALUES(?);

-- name: add-dirs
INSERT INTO dirs (id, parent_id, name, path, mod_time) VALUES

-- name: add-files
INSERT INTO files (id, parent_id, name, path, size, mod_time) VALUES

-- name: get-dirs-count
SELECT count(id) FROM dirs;

-- name: get-files-count
SELECT count(id) FROM files;

-- name: get-all-dirs
SELECT id, parent_id, name, path FROM dirs;

-- name: get-files-no-sha1
SELECT id, parent_id, name, path FROM files WHERE LENGTH(sha1) <= 0;

-- name: mod-file-sha1
UPDATE files SET sha1 = ? WHERE id = ?;

-- name: mod-file-status
UPDATE files SET status = ? WHERE id = ?;

-- name: trim-dir-ids
UPDATE dirs SET id = REPLACE(id, '-00000000', '-'), parent_id = REPLACE(parent_id, '-00000000', '-');

-- name: trim-file-ids
UPDATE files SET id = REPLACE(id, '-00000000', '-'), parent_id = REPLACE(parent_id, '-00000000', '-');

-- name: get-root-dir
SELECT id, parent_id, name, path FROM dirs WHERE parent_id = '0' LIMIT 1;

-- name: mod-root-dir
UPDATE dirs SET name = ?, path = ? WHERE parent_id = '0';

-- name: replace-dir-paths
UPDATE dirs SET path=REPLACE(path, ?, ?);

-- name: replace-file-paths
UPDATE files SET path=REPLACE(path, ?, ?);

-- name: check-table-exists
SELECT name FROM sqlite_master WHERE type='table' AND name=?;

-- name: get-db-version
SELECT db_version FROM infos LIMIT 1;

-- name: mod-db-version
UPDATE infos SET db_version=?;

-- name: migrate-v2
ALTER TABLE dirs
ADD COLUMN
"status" INTEGER NOT NULL DEFAULT 0;
