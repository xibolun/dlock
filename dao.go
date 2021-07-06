package dlock

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var r *Repo
var once sync.Once

// Repo mysql repo
type Repo struct {
	db *sql.DB
}

// LockTable table of lock
type LockTable struct {
	ID           int64
	Name         string
	LockResource string
	Host         string
	// unix nano time; since 1970-01-01 00:00:00
	ExpiredTime int64
	CreateAt    time.Time
	DeleteAt    *time.Time
}

const (
	insertSql = "INSERT INTO dlock (name, lock_resource, host , expire_at,created_at,deleted_at) VALUES (?, ?, ?, ?, ?,null)"
	querySql  = "select id, name, lock_resource,host ,expire_at,timestamp(created_at),deleted_at from dlock where name = ? and expire_at > ? for update "
	updateSql = "update dlock set deleted_at = ?  where id  =?"
	createSql = `
		create table dlock
		(
			id int(11) unsigned auto_increment comment '主键'
				primary key,
			created_at timestamp null comment '记录创建时间',
			deleted_at timestamp null comment '记录删除时间',
			name varchar(64) null comment '资源名称， lock key',
			lock_resource varchar(64) null comment '资源信息，lock value, uuid/code/......',
			host varchar(64) null comment '运行的主机,hostname or hostIp',
			expire_at int(11) null comment '过期时间'
		) comment '分布式锁' ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;`
)

// initRepo init database connection
// dsn  mysql dataSourceName
func initRepo(user, password, database, ip string, port int64) (*Repo, error) {
	if r != nil {
		return r, nil
	}

	Info("start to init repo.")

	var err error
	// init repo
	// Opening a driver typically will not attempt to connect to the database.
	dsn := assemblyDSN(user, password, ip, database, port)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(25)
	db.SetMaxOpenConns(25)

	// ping test
	if err := db.Ping(); err != nil {
		return nil, err
	}
	Info("ping database successful")

	once.Do(func() {
		//init table
		r = &Repo{db: db}
		if err = r.initTable(); err != nil {
			return
		}
	})
	return r, err
}

// assemblyDSN
// The default internal output type of MySQL DATE and DATETIME values is []byte
// which allows you to scan the value into a []byte, string or sql.RawBytes variable in your program.

// However, many want to scan MySQL DATE and DATETIME values into time.Time variables,
// which is the logical equivalent in Go to DATE and DATETIME in MySQL.
// You can do that by changing the internal output type from []byte to time.Time with the DSN parameter parseTime=true.
// You can set the default time.Time location with the loc DSN parameter.
// https://github.com/go-sql-driver/mysql#timetime-support
func assemblyDSN(dbUser, dbPassword, dbHost, dbName string, dbPort int64) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
}

// initTable init table
func (r *Repo) initTable() error {
	// check table is exist
	if exist, _ := r.checkTableIsNotExist(); exist {
		if err := r.createTable(); err != nil {
			return err
		}
	}

	return nil
}

// createTable check table is exist
func (r *Repo) createTable() error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(createSql)

	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer tx.Commit()
	return nil
}

// checkTableIsNotExist check table is exist
func (r *Repo) checkTableIsNotExist() (bool, error) {
	_, err := r.db.Query("select * from dlock")
	if err == nil {
		return false, nil
	}

	if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr != nil && mysqlErr.Number == 1146 {
		return true, nil
	}

	return false, err
}

// QueryLockRes
func (r *Repo) queryLockRes(cond *LockTable) (table *LockTable, err error) {
	tx, err := r.db.Begin()
	defer tx.Commit()
	if err != nil {
		return
	}

	rows, err := tx.Query(querySql, cond.Name, time.Now())
	if err != nil {
		_ = tx.Rollback()
		return
	}

	table = &LockTable{}
	for rows.Next() {
		if err = rows.Scan(&table.ID, &table.Name, &table.LockResource, &table.Host, &table.ExpiredTime, &table.CreateAt, &table.DeleteAt); err != nil {
			_ = tx.Rollback()
			_ = rows.Close()
			return
		}
	}

	return
}

// insertLockRes
func (r *Repo) insertLockRes(tab *LockTable) (affected int64, err error) {
	tx, err := r.db.Begin()
	defer tx.Commit()

	if err != nil {
		return
	}

	// check current_time timestamp after lock expire_time timestamp
	rows, err := tx.Query(querySql, tab.Name, time.Now().Unix())
	if err != nil {
		_ = tx.Rollback()
		return
	}

	table := &LockTable{}
	for rows.Next() {
		if err = rows.Scan(&table.ID, &table.Name, &table.LockResource, &table.Host, &table.ExpiredTime, &table.CreateAt, &table.DeleteAt); err != nil {
			_ = tx.Rollback()
			_ = rows.Close()
			return
		}
	}

	if table.ID > 0 {
		return 0, fmt.Errorf("%s is already exists", tab.Name)
	}

	result, err := tx.Exec(insertSql, tab.Name, tab.LockResource, tab.Host, tab.ExpiredTime, time.Now())
	if err != nil {
		_ = tx.Rollback()
		return
	}

	return result.LastInsertId()
}

// deleteLockRes
func (r *Repo) deleteLockRes(id int64) (affected int64, err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return
	}

	result, err := tx.Exec(updateSql, time.Now(), id)
	if err != nil {
		_ = tx.Rollback()
		return
	}

	return result.RowsAffected()
}

// deleteLockRes
func (r *Repo) deleteLockKey(key interface{}) (affected int64, err error) {
	tx, err := r.db.Begin()
	defer tx.Commit()

	if err != nil {
		return
	}

	// check current_time timestamp after lock expire_time timestamp
	rows, err := tx.Query(querySql, key, time.Now().Unix())
	if err != nil {
		_ = tx.Rollback()
		return
	}

	table := &LockTable{}
	for rows.Next() {
		if err = rows.Scan(&table.ID, &table.LockResource, &table.Host, &table.ExpiredTime, &table.CreateAt, &table.DeleteAt); err != nil {
			_ = tx.Rollback()
			_ = rows.Close()
			return
		}
	}

	if table.ID > 0 {
		return 0, nil
	}

	result, err := tx.Exec(updateSql, table.ID)
	if err != nil {
		_ = tx.Rollback()
		return
	}

	return result.LastInsertId()
}
