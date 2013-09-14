package etcdplus

import (
	"github.com/coreos/go-etcd/etcd"
	"strconv"
)

type Counter struct {
	client *etcd.Client
	key    string
}

func NewCounter(client *etcd.Client) *Counter {
	counter := Counter{
		client: client,
		key:    getUUID(),
	}

	client.Set(counter.key, "0", 0)
	return &counter
}

func (c *Counter) Increment() error {
	v, err := c.GetValue()
	if err != nil {
		return err
	}

	err = c.SetValue(v + 1)
	if err != nil {
		return err
	}

	return nil
}

func (c *Counter) Decrement() error {
	v, err := c.GetValue()
	if err != nil {
		return err
	}

	err = c.SetValue(v - 1)
	if err != nil {
		return err
	}

	return nil
}

func (c *Counter) GetValue() (int, error) {
	resp, err := c.client.Get(c.key)
	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(resp[0].Value)
	if err != nil {
		panic(err)
	}

	return i, nil
}

func (c *Counter) SetValue(value int) error {
	_, err := c.client.Set(c.key, strconv.Itoa(value), 0)
	if err != nil {
		return err
	}
	return nil
}
