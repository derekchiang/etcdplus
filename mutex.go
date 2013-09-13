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
		key:    uuid(),
	}

	client.Set(mutex.key, UNLOCK_VALUE, 0)
	return &mutex
}

func (c *Mutex) lock(timeout time.Duration) error {
	timeoutChan := make(<-chan time.Time)
	if timeout != 0 {
		timeoutChan = time.After(timeout)
	}

	for {
		select {
		case <-timeoutChan:
			return errors.New("Timeout")
		default:
			_, success, _ := c.client.TestAndSet(c.key,
				UNLOCK_VALUE, LOCK_VALUE, 0)
			if success {
				return nil
			}
		}
	}
}

func (c *Mutex) unlock() error {
	_, err := c.client.Set(c.key, UNLOCK_VALUE, 0)

	// TODO: should we keep trying in case err is not nil?
	return err
}
