package dlock

import (
	"testing"
	"time"
)

const (
	caFile   = ""
	certFile = ""
	keyFile  = ""
)

var endpoints = []string{"127.0.0.1:2379"}

func TestElock_Acquire(t *testing.T) {
	l, err := NewDLock(
		WithEtcdOption(dialTimeout, endpoints...),
		WithEtcdAuthOption("", "", caFile, certFile, keyFile, true))

	if err != nil {
		t.Error(err)
		return
	}

	success, err := l.Acquire(5*time.Minute, key, value, host)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("lock status : %t, value :%s", success, l.GetValue(key))

}

func TestELock_IsLock(t *testing.T) {
	l, _ := NewDLock(
		WithEtcdOption(dialTimeout, endpoints...),
		WithEtcdAuthOption("", "", caFile, certFile, keyFile, true))

	islock, err := l.IsLock(key)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("lock status : %t, value :%s", islock, l.GetValue(key))
}

func TestELock_Release(t *testing.T) {
	l, err := NewDLock(
		WithEtcdOption(dialTimeout, endpoints...),
		WithEtcdAuthOption("", "", caFile, certFile, keyFile, true))

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
