module libvirt-bosh-cpi

go 1.12

require (
	github.com/bmatcuk/doublestar v1.1.5 // indirect
	github.com/charlievieth/fs v0.0.0-20170613215519-7dc373669fa1 // indirect
	github.com/cloudfoundry/bosh-utils v0.0.0-20190803100152-d286f594c8d9
	github.com/cppforlife/bosh-cpi-go v0.0.0-20180718174221-526823bbeafd
	github.com/digitalocean/go-libvirt v0.0.0-20190626172931-4d226dd6c437 // GOOD
	//github.com/digitalocean/go-libvirt v0.0.0-20190715144809-7b622097a793 // BAD
	github.com/diskfs/go-diskfs v0.0.0-20190517155712-1190dcf1ff31
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d // indirect
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.3.0 // indirect
)

replace github.com/diskfs/go-diskfs => github.com/a2geek/go-diskfs v0.0.0-20190810191223-f09edeb3e6a4
