package db

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/gohutool/boot4go-util/db"
	_ "modernc.org/sqlite"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : _db.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/12 21:23
* 修改历史 : 1. [2022/5/12 21:23] 创建文件 by LongYong
*/

var dbPlus db.DBPlus

func init() {
	dbPath := "./config/data.db"
	if err := ensureDBFile(dbPath); err != nil {
		panic(err)
	}

	db1, err := sql.Open("sqlite", dbPath)
	if err != nil {
		panic(err)
	}
	dbPlus = db.DBPlus{DB: db1}
}

func ensureDBFile(dbPath string) error {
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0o755); err != nil {
		return err
	}

	if _, err := os.Stat(dbPath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	f, err := os.OpenFile(dbPath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	return f.Close()
}

var sql_table = `CREATE TABLE if not exists "t_user" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "userid" VARCHAR(64) NULL,
    "username" VARCHAR(64),
	"password" VARCHAR(128),
    "createtime" TIMESTAMP default (datetime('now', 'localtime'))
);

CREATE TABLE if not exists "t_db" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "dbid" VARCHAR(64) NULL,
    "endpoint" VARCHAR(64),
	"username" VARCHAR(64),
	"password" VARCHAR(64),
    "createtime" TIMESTAMP default (datetime('now', 'localtime'))
);

CREATE TABLE if not exists "t_repos" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "reposid" VARCHAR(64) NULL,
    "name" VARCHAR(255),
    "description" text,
    "endpoint" VARCHAR(255),
	"username" VARCHAR(255),
	"password" VARCHAR(255),
    "createtime" TIMESTAMP default (datetime('now', 'localtime'))
);

CREATE TABLE if not exists "t_orchestrator" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "orchestratorid" VARCHAR(64) NULL,
    "name" VARCHAR(64),
    "description" text,
    "json" text,
    "userid" VARCHAR(64),
    "createtime" TIMESTAMP default (datetime('now', 'localtime'))
);
`

func InitDB() {
	var userTableExists int
	if err := dbPlus.GetDB().QueryRow("select count(1) from sqlite_master where type='table' and name='t_user'").Scan(&userTableExists); err != nil {
		panic(err)
	}

	_, err := dbPlus.GetDB().Exec(sql_table)
	if err != nil {
		panic(err)
	}

	// Only create default admin user on first initialization (when schema is created).
	// If user records are deleted later, they won't be auto-restored on restart.
	if userTableExists == 0 {
		InitAdminUser()
	}
}
