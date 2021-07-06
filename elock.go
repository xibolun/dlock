package dlock

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
)

// eLock etcd lock
type eLock struct {
	// etcd cluster client
	ec  *clientv3.Client
	mux *sync.Mutex
	es  *eSession
}

type eSession struct {
	conMux *concurrency.Mutex
	s      *concurrency.Session
	ctx    context.Context
	k, v   string
	ttl    time.Duration
}

// global redis cluster client
var ec *clientv3.Client

// NewRLock create redis distributed lock
// options: other parameter configs
func NewELock(opts Options) (*eLock, error) {
	// require check
	if err := NewValidate().
		SliceEmpty(opts.Cluster, "etcd cluster").
		ToError();
		err != nil {
		return nil, err
	}

	var err error
	if ec == nil {
		once.Do(func() {
			if ec, err = clientv3.New(clientv3.Config{
				Endpoints:   opts.Cluster,
				DialTimeout: opts.DialTimeout,
				TLS:         nil,
				Username:    opts.User,
				Password:    opts.Password,
			}); err != nil {
				Errorf("init etcd client fail, %s", err.Error())
			}
		})

		// get member list,  test is connection success
		if ec == nil {
			return nil, fmt.Errorf("init etcd client fail, %s", err.Error())
		}
		if _, err := ec.MemberList(context.TODO()); err != nil {
			Errorf("etcd client conn fail, %s", err.Error())
			return nil, err
		}
	}

	return &eLock{
		ec:  ec,
		mux: &sync.Mutex{},
	}, nil
}

// Acquire 获取锁
func (l *eLock) Acquire(expiration time.Duration, key, value, node string) (bool, error) {

	var timeout context.Context
	var cancelFunc context.CancelFunc
	if expiration > 0 {
		timeout, cancelFunc = context.WithTimeout(context.Background(), expiration)
	} else {
		timeout = context.Background()
	}
	if cancelFunc != nil {
		defer cancelFunc()
	}

	// add lease
	rsp, err := l.ec.Grant(timeout, int64(expiration.Seconds()))
	if err != nil {
		return false, nil
	}

	s, err := concurrency.NewSession(l.ec, concurrency.WithLease(rsp.ID))
	if err != nil {
		return false, err
	}
	if s == nil {
		return false, err
	}

	m := concurrency.NewMutex(s, key)

	if err := m.Lock(timeout); err != nil {
		return false, err
	}

	l.mux.Lock()
	l.es = &eSession{
		conMux: m,
		s:      s,
		ctx:    timeout,
		k:      key,
		v:      value,
		ttl:    expiration,
	}
	l.mux.Unlock()

	return true, nil
}

// IsLock check if is locked already
// If key expire, redis will return ---> redis: nil
func (l *eLock) IsLock(key string) (bool, error) {
	return l.es != nil, nil
}

// UnLock release lock
func (l *eLock) UnLock(key string) (err error) {
	l.mux.Lock()
	if err = l.es.conMux.Unlock(l.es.ctx); err != nil {
		return
	}
	if err = l.es.s.Close(); err != nil {
		return
	}
	l.es.ctx.Done()
	l.es = nil
	l.mux.Unlock()

	return nil
}

// GetLockID  get lock value
func (l *eLock) GetValue(key string) (value string) {
	if l.es != nil && l.es.k == key {
		return l.es.v
	}
	return ""
}

// GetType  get lock type
func (l *eLock) GetType() string {
	return EtcdLockType
}
