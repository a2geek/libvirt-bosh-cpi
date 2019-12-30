package throttle

func NewNoneThrottle(config ThrottleConfig) Throttle {
	return noneThrottle{}
}

type noneThrottle struct {
}

func (n noneThrottle) Lock() error {
	// do nothing
	return nil
}

func (n noneThrottle) Unlock() error {
	// do nothing
	return nil
}
