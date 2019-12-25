package throttle

func NewNoopThrottle(config ThrottleConfig) Throttle {
	return noopThrottle{}
}

type noopThrottle struct {
}

func (n noopThrottle) Lock() error {
	// do nothing
	return nil
}

func (n noopThrottle) Unlock() error {
	// do nothing
	return nil
}
