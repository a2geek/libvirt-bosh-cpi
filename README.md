# libvirt-bosh-cpi

[![Go Report Card](https://goreportcard.com/badge/github.com/a2geek/libvirt-bosh-cpi)](https://goreportcard.com/report/github.com/a2geek/libvirt-bosh-cpi)
[![GitHub release](https://img.shields.io/github/v/release/a2geek/libvirt-bosh-cpi)](https://github.com/a2geek/libvirt-bosh-cpi/releases/latest)

A Go BOSH CPI for the [libvirt virtualization API](https://libvirt.org/).

# Motivation

Cloud Foundry is a very interesting tool and environment to target and use for development purposes. Cloud Foundry requires a BOSH CPI to deploy. Unfortunately, the existing BOSH CPI's seem to be development-only (that is, temporary), unsuitable for Cloud Foundry, or require a small cluster of machines, or deploy to a cloud provider.

Libvirt is a good compromise. The current intent is to keep it scoped at 1 host. If something larger exists, there are a number of solutions.

# Features

* Libvirt settings are configurable. Copy `manifests/libvirt-kvm-vars.yml` or `manifests/libvirt-qemu-vars.yml` to setup custom configurations.  Extend or alter as needed.
* Libvirt connectivity is configurable. Options include TCP or TLS connections. Unix sockets are available, but without being able to map the socket into the BOSH VM is unrealistic. See `manifests/libvirt-socket.yml`, `manifests/libvirt-tcp.yml`, and `manifests/libvirt-tls.yml`.
* Throttling capabilities. Slower disks become problematic for larger deployments. A simple file-locking mechanism throttles VM creation to not swamp slower hosts. See `manifests/throttle-file-lock.yml`. This capability is expandable if other needs arise.
* Ability to experiment with various stemcells. There are 3 stemcells which may work with Libvirt. See [Stemcell experiments](docs/STEMCELLS.md) for details.

# Libvirt Versions

| Distribution | Libvirt version | Notes |
| --- | --- | --- |
| Ubuntu 18.04 | 4.0.0 | Initial development version |
| Ubuntu 20.04 | 6.0.0 | Current development version |

# Stemcell Versions

| Stemcell version | Stemcell OS | Notes |
| --- | --- | --- |
| `621.x` | `ubuntu-xenial` | |
| `1.x` | `ubuntu-bionic` | Libvirt XML requires a graphics entry. See the `<graphics...>` entry in `libvirt_vm_domain` in [libvirt-kvm-vars.yml](manifests/libvirt-kvm-vars.yml). |

# Status

The CPI is capable of running and managing BOSH deployments, generally without issue. Longer term management (snapshotting vms, for instance) either are untested or unimplemented at this point.

Feel free to try it out! Feedback and PR's are welcome.

Notes:

> As of [v3](https://github.com/a2geek/libvirt-bosh-cpi/releases/tag/v3), disk management is configurable with throttling. This was done to help with deployments like Cloud Foundry that create 15 (or more!) VM's at once. With slower disks (non-NVMe) the machine may not be able to keep up. This keeps it manageable at the cost of slowing down VM creation.

> As of [v4](https://github.com/a2geek/libvirt-bosh-cpi/releases/tag/v4), disk assignment allows for disk resizing. Disks as still assigned statically; current assignment scheme:
>  * `/dev/vda`: boot disk
>  * `/dev/vdb`: ephemeral disk
>  * `/dev/vdc`: config disk
>  * `/dev/vdd` or `/dev/vde`: persistent disk (optional); two expected when resizing disks

# Tickets being watched

* Libvirt command switcheroo (work-around in place): [ticket](https://github.com/digitalocean/go-libvirt/issues/87)

## Closed

* Blocking IaaS provided disk resize: [ticket](https://github.com/cloudfoundry/bosh-agent/issues/221)
* Libvirt enhancement request to allow network ingress to be specified via API: [ticket](https://bugzilla.redhat.com/show_bug.cgi?id=1761123)
* Issues with TLS connections after a certain commit: [ticket](https://github.com/digitalocean/go-libvirt/issues/89)

# Documentation

* [Libvirt Configuration](docs/CONFIG.md)
* [Setup the BOSH Director](docs/INSTALL.md)
* [Useful Utilities](docs/UTILITIES.md)
* [Deploying Software to BOSH](docs/DEPLOYMENT.md)
* [Development notes](docs/DEVELOPING.md)
* [Stemcell experiments](docs/STEMCELLS.md)
* [TODO list](docs/TODO.md)
