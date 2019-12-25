package throttle

import (
	"context"
	"time"

	"github.com/gofrs/flock"
)

func NewFileLockThrottle(config ThrottleConfig) Throttle {
	return fileLockThrottle{
		fileLock: flock.New("/var/lock/libvirt-storage-volume.lock"),
	}
}

type fileLockThrottle struct {
	fileLock *flock.Flock
}

func (t fileLockThrottle) Lock() error {
	_, err := t.fileLock.TryLockContext(context.TODO(), 30*time.Second)
	return err
}

func (t fileLockThrottle) Unlock() error {
	return t.fileLock.Unlock()
}
