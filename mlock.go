package dlock

import (
	"fmt"
	"sync"
	"time"
)

// mLock  distributed lock
type mLock struct {
	id         int64
	namespace  string
	key        string
	host       string
	expireTime time.Duration
	repo       *Repo
	mux        *sync.RWMutex
}

// NewMLock create mysql distributed lock
// options: other parameter configs
func NewMLock(opts Options) (*mLock, error) {
	// require check
	if err := NewValidate().
		StringIsNull(opts.User, "database user").
		StringIsNull(opts.IP, "database host ip").
		StringIsNull(opts.Password, "database password").
		StringIsNull(opts.Name, "database value").ToError();
		err != nil {
		return nil, err
	}

	r, err := initRepo(opts.User, opts.Password, opts.Name, opts.IP, opts.Port)
	if err != nil {
		return nil, fmt.Errorf("init repo fail, err: %v", err)
	}

	if r == nil {
		return nil, fmt.Errorf("init repo fail")
	}

	return &mLock{
		repo: r,
		mux:  &sync.RWMutex{},
	}, nil
}

// Acquire 获取锁
func (l *mLock) Acquire(expiredTime time.Duration, key, value, host string) (bool, error) {
	id, err := l.repo.insertLockRes(&LockTable{Name: key, LockResource: value, ExpiredTime: time.Now().Add(expiredTime).Unix(), Host: host})
	if id > 0 {
		l.addLockID(id)
	}
	return id > 0 && err == nil, err
}

// IsLock check if is locked already
func (l *mLock) IsLock(key string) (bool, error) {
	tab, err := l.repo.queryLockRes(&LockTable{LockResource: key})
	return err == nil && tab != nil && tab.ID > 0, err
}

// UnLock release lock
func (l *mLock) UnLock(key string) error {
	_, err := l.repo.deleteLockKey(key)
	return err
}

// GetLockID get  lock id
func (l *mLock) GetValue(key string) (value string) {
	lock, _ := l.repo.queryLockRes(&LockTable{Name: key})
	if lock != nil {
		return lock.LockResource
	}
	return ""
}

// GetType  get lock type
func (l *mLock) GetType() string {
	return MysqlLockType
}

// addLockID 写入lock id
func (l *mLock) addLockID(id int64) {
	l.mux.Lock()
	l.id = id
	l.mux.Unlock()
}
