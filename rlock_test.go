package dlock

import (
	"testing"
	"time"
)

const (
	dialTimeout = 120
	cluster     = "10.0.3.59:7000,10.0.3.59:7001,10.0.3.59:7002,10.0.3.59:7003,10.0.3.59:7004,10.0.3.59:7005,10.0.3.59:7006"
)

func TestRLock_Acquire(t *testing.T) {
	l, err := NewDLock(
		WithRedisOption(password, dialTimeout, cluster))

	if err != nil {
		t.Error(err)
		return
	}

	success, err := l.Acquire(5*time.Minute, key, value, host)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("lock status : %t, value :%s", success, l.GetValue(value))
}

func TestRLock_Release(t *testing.T) {
	l, err := NewDLock(
		WithRedisOption(password, dialTimeout, cluster))

	if err != nil {
		t.Error(err)
		return
	}

	success, err := l.Acquire(5*time.Minute, key, value, host)
	if err != nil {
		t.Error(err)
		return
	}
	if success {
		time.Sleep(30 * time.Second)
		if err = l.UnLock(value); err != nil {
			t.Error(err)
			return
		}
	}
}
