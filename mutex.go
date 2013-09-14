package etcdplus

import (
	"errors"
	// "fmt"
	"github.com/coreos/go-etcd/etcd"
	"time"
)

const (
	UNLOCK_VALUE = "unlocked"
)

type Mutex struct {
	// key is a well-known string that nodes can use to
	// access the mutex.
	// lockVal is specific to a node.  When a node locks
	// a mutex, the value is set to lockVal; since lockVal
	// is specific to that node, other nodes won't be able
	// to unlock the mutex.
	client  *etcd.Client
	key     string
	lockVal string
}

func NewMutex(client *etcd.Client, key string) *Mutex {
	// Generate a random key if no key is given.
	// The user should inspect the returned mutex
	// and tell other nodes about the key.
	if key == "" {
		key = getUUID()
	}

	mutex := Mutex{
		client:  client,
		key:     key,
		lockVal: getUUID(),
	}

	return &mutex
}

func (m *Mutex) Lock(timeout time.Duration) error {
	timeoutChan := getTimeoutChan(timeout)
	successChan := make(chan bool)
	stopChan := make(chan bool) // Stop the goroutine

	var lockHelper func()
	lockHelper = func() {
		resps, err := m.client.Get(m.key)
		if err != nil {
			etcdErr := err.(etcd.EtcdError)
			if etcdErr.ErrorCode == 100 {
				// Create the mutex.
				// A mutex might be missing for two reasons:
				// 1. This is the first time anyone ever tries to
				// lock the mutex.
				// 2. Some node died when holding the mutex.
				m.client.TestAndSet(m.key, "", UNLOCK_VALUE, 0)

				// Restart
				lockHelper()
				return
			}
			panic(err)
		}

		resp := resps[0]
		for {
			sinceIndex := resp.Index + 1
			// Try to lock if it's unlocked
			if resp.Value == UNLOCK_VALUE {
				_, success, _ := m.client.TestAndSet(m.key,
					UNLOCK_VALUE, m.lockVal, 0)
				if success {
					successChan <- true
					return
				}
			}

			// Wait till it's changed
			resp, err = m.client.Watch(m.key, sinceIndex, nil, stopChan)
			if err != nil {
				// The reason is likely that stopChan was used
				return
			}
		}
	}

	go lockHelper()

	select {
	case <-timeoutChan:
		stopChan <- true
		return errors.New("Timeout")
	case <-successChan:
		return nil
	}
}

func (m *Mutex) Unlock() error {
	_, _, err := m.client.TestAndSet(m.key, m.lockVal, UNLOCK_VALUE, 0)

	// TODO: should we keep trying in case err is not nil?
	return err
}
