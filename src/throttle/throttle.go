package throttle

const ThrottleNone = "none"
const ThrottleFileLock = "file-lock"

type Throttle interface {
	Lock() error
	Unlock() error
}

type ThrottleConfig struct {
	Name string
}

func NewThrottle(config ThrottleConfig) Throttle {
	switch config.Name {
	case ThrottleNone:
		return NewNoopThrottle(config)
	case ThrottleFileLock:
		return NewFileLockThrottle(config)
	default:
		return NewNoopThrottle(config)
	}
}
