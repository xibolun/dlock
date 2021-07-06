package dlock

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v7"
)

// DLock distributed lock interface
type DLock interface {
	// Acquire:  get a lock, if success return true
	// expiration required
	// key: lock resource name , required
	// value: lock resource, required
	// host: which host need this lock resource, omit
	// namespace: which app need this lock resource, omit
	Acquire(expiration time.Duration, key, value, host string) (bool, error)
	IsLock(key string) (bool, error)
	UnLock(key string) error
	GetValue(key string) string
	GetType() string
}

// dlock  distributed lock
type dlock struct {
	// MysqlLockType return id
	// RedisLockType hash value
	id        int64
	namespace string
	// may be interface
	lockResource interface{}
	hostIP       string
	// timestamp
	expireTime int64

	repo *Repo

	// redis cluster client
	rc  *redis.ClusterClient
	mux *sync.RWMutex
}

const (
	MysqlLockType = "mysql"
	RedisLockType = "redis"
	EtcdLockType  = "etcd"
)

var (
	NotSupportedTypeLockErr = fmt.Errorf("not support this type distibuted lock")
)

// NewDLock create distributed lock
// options: other parameter configs
func NewDLock(options ...func(*Options)) (DLock, error) {
	// init database
	var opts Options
	for i := range options {
		options[i](&opts)
	}

	var dlock DLock
	var err error
	switch opts.Type {
	case MysqlLockType:
		dlock, err = NewMLock(opts)
	case EtcdLockType:
		dlock, err = NewELock(opts)
	default:
		dlock, err = NewRLock(opts)
	}

	return dlock, err
}
