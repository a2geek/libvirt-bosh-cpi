# libvirt-bosh-cpi
A Go BOSH CPI for the [libvirt virtualization API](https://libvirt.org/).

## Setup

Since the packaging is a bit wonky in this setup (call it _vendoring_), the `GOPATH` needs to be set to the project root. In VS Code, the `settings.json` file will contain something like:

```
{
    "go.gopath": "/path/to/the/directory/libvirt-bosh-cpi"
}
```
