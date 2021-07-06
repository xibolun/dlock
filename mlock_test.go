package dlock

import (
	"testing"
	"time"
)

const (
	key   = "job_id"
	value = "033d8776c308402a9b20a4dd970f8742"
	host  = "10.0.3.57"
)

func TestMLock_Acquire(t *testing.T) {
	l, err := NewDLock(
		WithDBOption(user, password, ip, database, port))

	if err != nil {
		t.Error(err)
		return
	}

	success, err := l.Acquire(5*time.Minute, key, value, host)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("lock status : %t, value:%s", success, l.GetValue(value))
}

func TestMLock_Release(t *testing.T) {
	l, err := NewDLock(
		WithDBOption(user, password, ip, database, port))

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
