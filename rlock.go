package dlock

import (
	"sync"
	"time"

	"github.com/go-redis/redis/v7"
)

// rLock redis lock
type rLock struct {
	// redis cluster client
	rc  Clienter
	mux *sync.Mutex

	key        string
	value      interface{}
	expiration time.Duration
}

// global redis cluster client
var rc Clienter

// NewRLock create redis distributed lock
// options: other parameter configs
func NewRLock(opts Options) (*rLock, error) {
	// require check
	if err := NewValidate().
		SliceEmpty(opts.Cluster, "redis cluster").
		// StringIsNull(opts.Password, "redis password").
		ToError();
		err != nil {
		return nil, err
	}

	if rc == nil {
		once.Do(func() {
			if len(opts.Cluster) > 1 {
				rc = redis.NewClusterClient(&redis.ClusterOptions{
					Addrs:       opts.Cluster,
					Password:    opts.Password,
					DialTimeout: opts.DialTimeout,
				})
			} else {
				rc = redis.NewClient(&redis.Options{
					Addr:        opts.Cluster[0],
					Password:    opts.Password,
					DialTimeout: opts.DialTimeout,
				})
			}

		})

		// client ping
		// ignore result string PONG
		if _, err := rc.Ping().Result(); err != nil {
			Errorf("redis cluster client ping fail, %v", err)
			return nil, err
		}
	}

	return &rLock{
		rc:  rc,
		mux: &sync.Mutex{},
	}, nil
}

// Acquire 获取锁
func (l *rLock) Acquire(expiration time.Duration, key, value, host string) (bool, error) {
	//// TODO 添加自动续租功能
	succ, err := l.rc.SetNX(key, value, expiration).Result()
	if err != nil {
		l.rc.Del(l.key)
	}
	return succ, err
}

// IsLock check if is locked already
// If key expire, redis will return ---> redis: nil
func (l *rLock) IsLock(key string) (bool, error) {
	if str, err := l.rc.Get(key).Result(); err != nil {
		if err.Error() == "redis: nil" {
			return false, nil
		}
		return false, err
	} else {
		return len(str) > 0, nil
	}
}

// UnLock release lock
func (l *rLock) UnLock(key string) (err error) {
	return l.rc.Del(key).Err()
}

// GetLockID  get lock value
func (l *rLock) GetValue(key string) (value string) {
	return l.rc.Get(key).String()
}

// GetType  get lock type
func (l *rLock) GetType() string {
	return RedisLockType
}

// Clienter  redis client
// adapter ClusterClient && Client
type Clienter interface {
	SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	Del(keys ...string) *redis.IntCmd
	Get(key string) *redis.StringCmd
	Ping() *redis.StatusCmd
}
