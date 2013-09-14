package etcdplus

import (
	"errors"
	"github.com/coreos/go-etcd/etcd"
	"time"
)

const (
	LOCK_VALUE   = "lock"
	UNLOCK_VALUE = "unlock"
)

type Mutex struct {
	client *etcd.Client
	key    string
}

func NewMutex(client *etcd.Client) *Mutex {
	mutex := Mutex{
		client: client,
		key:    getUUID(),
	}

	client.Set(mutex.key, UNLOCK_VALUE, 0)
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
			panic(err)
		}

		if len(resps) != 1 {
			panic("Multiple responses")
		}

		resp := resps[0]
		for {
			sinceIndex := resp.Index + 1
			// Try to lock if it's unlocked
			if resp.Value == UNLOCK_VALUE {
				_, success, _ := m.client.TestAndSet(m.key,
					UNLOCK_VALUE, LOCK_VALUE, 0)
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
	_, err := m.client.Set(m.key, UNLOCK_VALUE, 0)

	// TODO: should we keep trying in case err is not nil?
	return err
}
