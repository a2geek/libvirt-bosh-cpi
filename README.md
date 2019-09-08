# libvirt-bosh-cpi

A Go BOSH CPI for the [libvirt virtualization API](https://libvirt.org/).

# Status

Many things do work at this point and the CPI is pretty stable. Most likely a BOSH deployment will successfully _deploy_. Longer term management (resizing disks or snapshotting vms, for instance) either are untested or unimplemented at this point.

## Documentation

* [Libvirt Configuration](docs/CONFIG.md)
* [Setup the BOSH Director](docs/INSTALL.md)
* [Useful Utilities](docs/UTILITIES.md)
* [Deploying Software to BOSH](docs/DEPLOYMENT.md)
* [Development notes](docs/DEVELOPING.md)
* [Stemcell experiments](docs/STEMCELLS.md)
* [TODO list](docs/TODO.md)
