package throttle

import "time"

const ThrottleNone = "none"
const ThrottleFileLock = "file-lock"

type Throttle interface {
	Lock() error
	Unlock() error
}

type ThrottleConfig struct {
	// Throttle config name: "none", "file-lock", all others are defaulted.
	Name string
	// Path for "file-lock". Default is "/var/lock/libvirt-storage-volume.lock".
	Path string
	// Retry is the duration for "file-lock" to re-attempt the lock. Default is "30*time.Second".
	Retry time.Duration
}

func NewThrottle(config ThrottleConfig) Throttle {
	switch config.Name {
	case ThrottleNone:
		return NewNoneThrottle(config)
	case ThrottleFileLock:
		return NewFileLockThrottle(config)
	default:
		return NewNoneThrottle(config)
	}
}
