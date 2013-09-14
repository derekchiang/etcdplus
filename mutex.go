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

func (m *Mutex) Unlock() error {
	_, err := m.client.Set(m.key, UNLOCK_VALUE, 0)

	// TODO: should we keep trying in case err is not nil?
	return err
}
