package throttle

import (
	"context"
	"time"

	"github.com/gofrs/flock"
)

func NewFileLockThrottle(config ThrottleConfig) Throttle {
	if config.Path == "" {
		config.Path = "/var/lock/libvirt-storage-volume.lock"
	}
	if config.Retry == 0 {
		config.Retry = 30 * time.Second
	}
	return fileLockThrottle{
		fileLock: flock.New(config.Path),
		retry:    config.Retry,
	}
}

type fileLockThrottle struct {
	fileLock *flock.Flock
	retry    time.Duration
}

func (t fileLockThrottle) Lock() error {
	_, err := t.fileLock.TryLockContext(context.TODO(), t.retry)
	return err
}

func (t fileLockThrottle) Unlock() error {
	return t.fileLock.Unlock()
}
