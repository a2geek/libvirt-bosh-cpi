# libvirt-bosh-cpi

A Go BOSH CPI for the [libvirt virtualization API](https://libvirt.org/).

# Motivation

Cloud Foundry is a very interesting tool and environment to target and use for development purposes. Cloud Foundry requires a BOSH CPI to deploy. Unfortunately, the existing BOSH CPI's seem to be development-only (that is, temporary), unsuitable for Cloud Foundry, or require a small cluster of machines, or deploy to a cloud provider.

Libvirt is a good compromise. The current intent is to keep it scoped at 1 host. If something larger exists, there are a number of solutions.

# Status

Many things do work at this point and the CPI is pretty stable. Most likely a BOSH deployment will successfully _deploy_. Longer term management (resizing disks or snapshotting vms, for instance) either are untested or unimplemented at this point.

> Note that as of [v2](https://github.com/a2geek/libvirt-bosh-cpi/releases/tag/v2), the CPI only manages one disk at a time. This was done to help with deployments like Cloud Foundry that create 15 (or more!) VM's at once. With slower disks (non-NVMe) the machine may not be able to keep up. This keeps it manageable at the cost of slowing down VM creation.

Feel free to try it out! Feedback and PR's are welcome.

# Tickets being watched

* Blocking IaaS provided disk resize: [ticket](https://github.com/cloudfoundry/bosh-agent/issues/221)
* Issues with TLS connections after a certain commit: [ticket](https://github.com/digitalocean/go-libvirt/issues/89)
* Libvirt command switcheroo (work-around in place): [ticket](https://github.com/digitalocean/go-libvirt/issues/87)

## Closed

* Libvirt enhancement request to allow network ingress to be specified via API: [ticket](https://bugzilla.redhat.com/show_bug.cgi?id=1761123)

# Documentation

* [Libvirt Configuration](docs/CONFIG.md)
* [Setup the BOSH Director](docs/INSTALL.md)
* [Useful Utilities](docs/UTILITIES.md)
* [Deploying Software to BOSH](docs/DEPLOYMENT.md)
* [Development notes](docs/DEVELOPING.md)
* [Stemcell experiments](docs/STEMCELLS.md)
* [TODO list](docs/TODO.md)
