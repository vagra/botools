-- name: ----------------------------------------
-- CREATE TABLES
-------------------------------------------------

-- name: create-dirs-table
BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS "dirs" (
	"id" 		VARCHAR(32) NOT NULL UNIQUE,
	"parent_id" VARCHAR(32) NOT NULL DEFAULT '0',
	"name" 		TEXT NOT NULL DEFAULT '',
	"path" 		TEXT NOT NULL DEFAULT '',
	"size" 		INTEGER NOT NULL DEFAULT 0,
	"status" 	INTEGER NOT NULL DEFAULT 0,
	"error" 	INTEGER NOT NULL DEFAULT 0,
	"mod_time"	VARCHAR(32) NOT NULL DEFAULT '',
	PRIMARY KEY("id")
);

COMMIT;

-- name: create-files-table
BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS "files" (
	"id" 		VARCHAR(32) NOT NULL UNIQUE,
	"parent_id" VARCHAR(32) NOT NULL DEFAULT '',
	"name" 		TEXT NOT NULL DEFAULT '',
	"path" 		TEXT NOT NULL DEFAULT '',
	"size" 		INTEGER NOT NULL DEFAULT 0,
	"status" 	INTEGER NOT NULL DEFAULT 0,
	"error" 	INTEGER NOT NULL DEFAULT 0,
	"dup_id" 	VARCHAR(32) NOT NULL DEFAULT '',
	"sha1" 		VARCHAR(64) NOT NULL DEFAULT '',
	"mod_time" 	VARCHAR(32) NOT NULL DEFAULT '',
	PRIMARY KEY("id")
);

COMMIT;

-- name: create-vdirs-table
BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS "vdirs" (
	"id"		VARCHAR(32) NOT NULL UNIQUE,
	"parent_id"	VARCHAR(32) NOT NULL DEFAULT '0',
	"name"		TEXT NOT NULL DEFAULT '',
	"path"		TEXT NOT NULL DEFAULT '',
	"status"	INTEGER NOT NULL DEFAULT 0,
	"mod_time"	VARCHAR(32) NOT NULL DEFAULT '',
	PRIMARY KEY("id")
);

COMMIT;

-- name: create-vfiles-table
BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS "vfiles" (
	"id"		VARCHAR(32) NOT NULL UNIQUE,
	"real_id"	VARCHAR(32) NOT NULL DEFAULT '',
	"parent_id"	VARCHAR(32) NOT NULL DEFAULT '',
	"name"		TEXT NOT NULL DEFAULT '',
	"path"		TEXT NOT NULL DEFAULT '',
	"status"	INTEGER NOT NULL DEFAULT 0,
	"mod_time"	VARCHAR(32) NOT NULL DEFAULT '',
	PRIMARY KEY("id")
);

COMMIT;

-- name: create-infos-table
BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS "infos" (
	"id" INTEGER	NOT NULL UNIQUE,
	"db_version"	INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY("id" AUTOINCREMENT)
);

INSERT INTO infos (db_version) VALUES(4);

COMMIT;

-- name: ----------------------------------------
-- COMMON
-------------------------------------------------

-- name: init-pragma
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;

-- name: begin
BEGIN

-- name: end
END

-- name: check-table-exists
SELECT name FROM sqlite_master WHERE type='table' AND name=?;

-- name: ----------------------------------------
-- INSERT
-------------------------------------------------

-- name: add-dir
INSERT INTO dirs (id, parent_id, name, path, mod_time) VALUES(?, ?, ?, ?, ?);

-- name: add-file
INSERT INTO files (id, parent_id, name, path, size, mod_time) VALUES(?, ?, ?, ?, ?, ?);

-- name: add-info
INSERT INTO infos (db_version) VALUES(?);

-- name: ----------------------------------------
-- BATCH INSERT
-------------------------------------------------

-- name: add-dirs
INSERT INTO dirs (id, parent_id, name, path, mod_time) VALUES

-- name: add-files
INSERT INTO files (id, parent_id, name, path, size, mod_time) VALUES

-- name: ----------------------------------------
-- QUERY
-------------------------------------------------

-- name: -------- query dirs --------

-- name: get-root-dir
SELECT id, parent_id, name, path, status, error FROM dirs WHERE parent_id = '0' LIMIT 1;

-- name: get-dirs-count
SELECT count(id) FROM dirs;

-- name: get-max-dir-id
SELECT MAX(id) FROM dirs;

-- name: get-all-dirs
SELECT id, parent_id, name, path, status, error FROM dirs;

-- name: get-dir-id-from-path
SELECT id FROM dirs WHERE path = ? LIMIT 1;

-- name: get-a-dir-id
SELECT id FROM dirs LIMIT 1;

-- name: get-next-dir
SELECT id, parent_id, name, path, status, error FROM dirs WHERE id > ? ORDER BY id LIMIT 1;

-- name: -------- query files --------

-- name: get-files-count
SELECT count(id) FROM files;

-- name: get-max-file-id
SELECT MAX(id) FROM files;

-- name: get-all-files
SELECT id, parent_id, name, path, status, error FROM files;

-- name: get-unique-files
SELECT id, parent_id, name, path, status, error FROM files WHERE dup_id = 'UNIQUE' AND status = 0 AND error = 0;

-- name: get-unique-or-error-files
SELECT id, parent_id, name, path, status, error FROM files WHERE (dup_id = 'UNIQUE' AND status = 0) OR error = 1;

-- name: get-no-sha1-files-count
SELECT count(id) FROM files WHERE LENGTH(sha1) <= 0 AND status = 0;

-- name: get-no-sha1-files
SELECT id, parent_id, name, path FROM files WHERE LENGTH(sha1) <= 0 AND status = 0;

-- name: get-file-id-from-path
SELECT id FROM files WHERE path = ? LIMIT 1;

-- name: get-a-file-id
SELECT id FROM files LIMIT 1;

-- name: get-next-file
SELECT id, parent_id, name, path, size, status, error, dup_id, sha1 FROM files WHERE id > ? ORDER BY id LIMIT 1;

-- name: get-next-nodup-file
SELECT id, parent_id, name, path, size, status, error, dup_id, sha1 FROM files WHERE (LENGTH(dup_id) < 4 AND LENGTH(sha1) > 8 AND status = 0 AND error = 0) AND id > ? ORDER BY id LIMIT 1;

-- name: get-next-dup-file
SELECT id, parent_id, name, path, size, status, error, dup_id, sha1 FROM files WHERE (LENGTH(dup_id) > 8 AND status = 0 AND error = 0) AND id > ? ORDER BY id LIMIT 1;

-- name: get-next-exist-or-error-file
SELECT id, parent_id, name, path, size, status, error, dup_id, sha1 FROM files WHERE (status = 0 or error = 1) AND id > ? ORDER BY id LIMIT 1;

-- name: -------- query infos --------

-- name: get-db-version
SELECT db_version FROM infos LIMIT 1;

-- name: ----------------------------------------
-- UPDATE
-------------------------------------------------

-- name: -------- update dirs --------

-- name: mod-root-dir
UPDATE dirs SET name = ?, path = ? WHERE parent_id = '0';

-- name: trim-dirs-id
UPDATE dirs SET id = REPLACE(id, '-00000000', '-'), parent_id = REPLACE(parent_id, '-00000000', '-');

-- name: mod-dirs-status
UPDATE dirs SET status = ?;

-- name: replace-dirs-path
UPDATE dirs SET path = ( ? || substr(path, length(?)+1) ) WHERE path LIKE (? || '%');

-- name: replace-dirs-id
UPDATE dirs SET id = REPLACE(id, ?, ?);

-- name: replace-dirs-parent-id
UPDATE dirs SET parent_id = REPLACE(parent_id, ?, ?);

-- name: mod-dir-status
UPDATE dirs SET status = ? WHERE id = ?;

-- name: mod-dir-error
UPDATE dirs SET error = ? WHERE id = ?;

-- name: -------- update files --------

-- name: trim-files-id
UPDATE files SET id = REPLACE(id, '-00000000', '-'), parent_id = REPLACE(parent_id, '-00000000', '-');

-- name: mod-files-status
UPDATE files SET status = ?;

-- name: reset-files-dup-id
UPDATE files SET dup_id = '';

-- name: replace-files-path
UPDATE files SET path = ( ? || substr(path, length(?)+1) ) WHERE path LIKE (? || '%');

-- name: replace-files-id
UPDATE files SET id = REPLACE(id, ?, ?);

-- name: replace-files-parent-id
UPDATE files SET parent_id = REPLACE(parent_id, ?, ?);

-- name: mod-dir-files-error
UPDATE files SET error = ? WHERE parent_id = ?;

-- name: mod-file-sha1
UPDATE files SET sha1 = ? WHERE id = ?;

-- name: mod-file-status
UPDATE files SET status = ? WHERE id = ?;

-- name: mod-file-error
UPDATE files SET error = ? WHERE id = ?;

-- name: mod-file-dup-id
UPDATE files SET dup_id = ? WHERE id = ?;

-- name: -------- update infos --------

-- name: mod-db-version
UPDATE infos SET db_version=?;

-- name: ----------------------------------------
-- MIGRATIONS
-------------------------------------------------

-- name: migrate-v2
BEGIN TRANSACTION;

ALTER TABLE dirs
ADD COLUMN "status" INTEGER NOT NULL DEFAULT 0;

UPDATE infos SET db_version=3;

COMMIT;

-- name: migrate-v3
BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS "vdirs" (
	"id"		VARCHAR(32) NOT NULL UNIQUE,
	"parent_id"	VARCHAR(32) NOT NULL DEFAULT '0',
	"name"		TEXT NOT NULL DEFAULT '',
	"path"		TEXT NOT NULL DEFAULT '',
	"status"	INTEGER NOT NULL DEFAULT 0,
	"mod_time"	VARCHAR(32) NOT NULL DEFAULT '',
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "vfiles" (
	"id"		VARCHAR(32) NOT NULL UNIQUE,
	"real_id"	VARCHAR(32) NOT NULL DEFAULT '',
	"parent_id"	VARCHAR(32) NOT NULL DEFAULT '',
	"name"		TEXT NOT NULL DEFAULT '',
	"path"		TEXT NOT NULL DEFAULT '',
	"status"	INTEGER NOT NULL DEFAULT 0,
	"mod_time"	VARCHAR(32) NOT NULL DEFAULT '',
	PRIMARY KEY("id")
);

UPDATE infos SET db_version=3;

COMMIT;

-- name: migrate-v4

BEGIN TRANSACTION;

ALTER TABLE dirs
ADD COLUMN "error" 	INTEGER NOT NULL DEFAULT 0;

ALTER TABLE files
ADD COLUMN "error" 	INTEGER NOT NULL DEFAULT 0;
ALTER TABLE files
ADD COLUMN "dup_id" VARCHAR(32) NOT NULL DEFAULT '';

UPDATE infos SET db_version=4;

COMMIT;