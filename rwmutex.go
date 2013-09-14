package etcdplus

import (
	"errors"
	"github.com/coreos/go-etcd/etcd"
	"time"
)

const (
	R_LOCK_VALUE = "rlock"
	W_LOCK_VALUE = "wlock"
)

// A classic readers-writer mutex that allows multiple readers
// or a single writer to hold the mutex.
type RWMutex Mutex

// Return a new RWMutex.
func NewRWMutex(client *etcd.Client) *RWMutex {
	mutex := RWMutex{
		client: client,
		key:    getUUID(),
	}

	client.Set(mutex.key, UNLOCK_VALUE, 0)
	return &mutex
}

func (m *RWMutex) RLock(timeout time.Duration) error {
	timeoutChan := getTimeoutChan(timeout)

	for {
		select {
		case <-timeoutChan:
			return errors.New("Timeout")
		default:
			_, success, _ := m.client.TestAndSet(m.key,
				UNLOCK_VALUE, LOCK_VALUE, 0)
			if success {
				return nil
			}
		}
	}
}
