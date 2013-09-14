package etcdplus

import (
	"github.com/coreos/go-etcd/etcd"
	"testing"
)

func TestCounter(t *testing.T) {
	c := etcd.NewClient()
	counter := NewCounter(c, "")

	// Init value should be 0
	v, err := counter.GetValue()
	if err != nil {
		t.Fatal(err)
	}

	if v != 0 {
		t.Fatal("Init value should be 0.")
	}

	err = counter.Increment()
	if err != nil {
		t.Fatal(err)
	}

	v, err = counter.GetValue()
	if err != nil {
		t.Fatal(err)
	}

	if v != 1 {
		t.Fatal("Should be 1.")
	}

	err = counter.Decrement()
	if err != nil {
		t.Fatal(err)
	}

	v, err = counter.GetValue()
	if err != nil {
		t.Fatal(err)
	}

	if v != 0 {
		t.Fatal("Should be 0.")
	}
}
