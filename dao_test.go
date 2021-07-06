package dlock

import (
	"fmt"
	"testing"
	"time"
)

const (
	user     = "root"
	password = "Yunjikeji#123"
	ip       = "10.0.2.8"
	database = "cloudboot_3.0.0"
	port     = 3306
)

func Test_initRepo(t *testing.T) {
	_, err := initRepo(user, password, database, ip, port)

	if err != nil {
		t.Error(err)
		return
	}
}

func Test_createTable(t *testing.T) {
	r, _ := initRepo(user, password, database, ip, port)

	err := r.createTable()
	if err != nil {
		t.Error(err)
		return
	}
}

func Test_checkTableIsNotExist(t *testing.T) {
	r, _ := initRepo(user, password, database, ip, port)

	exists, err := r.checkTableIsNotExist()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(exists)
}

func Test_queryLockRes(t *testing.T) {
	r, _ := initRepo(user, password, database, ip, port)

	table, err := r.queryLockRes(&LockTable{ LockResource: "38", ExpiredTime: time.Now().UnixNano()})
	if err != nil {
		t.Error(err)
		return
	}
	// output result
	fmt.Printf("id:%d\name:%s\nhost_ip:%s\nexpire_at:%d\ncreate_at:%s\ndelete_at:%v\n", table.ID, table.LockResource, table.Host, table.ExpiredTime, table.CreateAt, table.DeleteAt)
}

func Test_insertLockRes(t *testing.T) {
	r, _ := initRepo(user, password, database, ip, port)

	exists, err := r.insertLockRes(&LockTable{
		LockResource: "uuid",
		ExpiredTime:  time.Now().UnixNano(),
		Host:         "10.0.3.37",
	})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(exists)
}

func Test_deleteLockRes(t *testing.T) {
	r, _ := initRepo(user, password, database, ip, port)

	affected, err := r.deleteLockRes(1)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(affected)
}

// MysqlLockType s lock test
func Test_SLock(t *testing.T) {
	r, _ := initRepo(user, password, database, ip, port)

	tx, err := r.db.Begin()
	if err != nil {
		return
	}

	rows, err := tx.Query(querySql, "aa", "bb", "cc")
	if err != nil {
		_ = tx.Rollback()
		return
	}

	// sleep 10s
	// then exec "select * from xxx for update" for testing
	time.Sleep(10 * time.Second)
	fmt.Println(rows.Next())

	defer tx.Commit()

	return
}
