package etcdplus

import (
	"github.com/coreos/go-etcd/etcd"
	"testing"
	"time"
)

func TestMutex(t *testing.T) {
	c := etcd.NewClient()
	mutex := NewMutex(c)

	err := mutex.lock(0)
	if err != nil {
		t.Fatal(err)
	}

	fatal := make(chan bool)
	go func() {
		err := mutex.lock(1 * time.Second) // should not succeed
		if err == nil {
			fatal <- true
		} else {
			fatal <- false
		}
	}()

	if <-fatal {
		t.Fatal("Should not be able to acquire a lock twice.")
	}

	success := make(chan bool)
	go func() {
		err := mutex.lock(1)
		if err == nil {
			success <- true
		} else {
			success <- false
		}
	}()

	err = mutex.unlock()
	if err != nil {
		t.Fatal(err)
	}

	if !(<-success) {
		t.Fatal("Fail to acquire lock after it's released.")
	}
}