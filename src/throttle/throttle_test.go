package throttle_test

import (
	"fmt"
	"libvirt-bosh-cpi/throttle"
	"testing"
	"time"
)

const delay = 2 * time.Second
const samples = 5

func TestDefaultThrottle(t *testing.T) {
	fmt.Println("Testing default throttle. Expected duration approximately 2 seconds.")
	testThrottle(t, throttle.ThrottleConfig{Name: "default"}, func(d time.Duration) bool {
		return d >= delay && d <= 2*delay
	})
}

func TestFilelockThrottle(t *testing.T) {
	fmt.Println("Testing 'file-lock' throttle. Expected duration approximately 14 seconds.")
	config := throttle.ThrottleConfig{
		Name:  "file-lock",
		Path:  "/tmp/libvirt-storage-volume.lock",
		Retry: 2 * time.Second,
	}
	testThrottle(t, config, func(d time.Duration) bool {
		return d > delay*samples
	})
}

func TestNoneThrottle(t *testing.T) {
	fmt.Println("Testing 'none' throttle. Expected duration approximately 2 seconds.")
	testThrottle(t, throttle.ThrottleConfig{Name: "none"}, func(d time.Duration) bool {
		return d >= delay && d <= 2*delay
	})
}

func testThrottle(t *testing.T, config throttle.ThrottleConfig, validate func(d time.Duration) bool) {
	ch := make(chan bool)
	for i := 0; i < samples; i++ {
		go testOne(config, ch)
	}

	start := time.Now()
	actualSuccess := 0
	for i := 0; i < samples; i++ {
		b := <-ch
		if b {
			actualSuccess++
		}
	}

	duration := time.Since(start)
	if !validate(duration) {
		t.Errorf("Actual lock duration did not meet expected range")
	}
}

func testOne(config throttle.ThrottleConfig, r chan bool) {
	th := throttle.NewThrottle(config)
	defer th.Unlock()
	err := th.Lock()
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(delay)
	r <- err == nil
}
