package etcdplus

import (
	"github.com/nu7hatch/gouuid"
	"time"
)

// TODO
func getUUID() string {
	uid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return uid.String()
}

// Return a channel that will send back a value when
// timeout.  If timeout is 0, the value will never arrive.
func getTimeoutChan(timeout time.Duration) <-chan time.Time {
	timeoutChan := make(<-chan time.Time)
	if timeout != 0 {
		timeoutChan = time.After(timeout)
	}
	return timeoutChan
}
