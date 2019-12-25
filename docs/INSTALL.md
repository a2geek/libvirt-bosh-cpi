# Installation

This is all setup for Ubuntu Bionic (18.04) using QEMU/KVM for the Libvirt virtualization component and Openstack stemcell.

All commands should be run from the root of this repository.

1. Set the location of `BOSH_DEPLOYMENT_DIR` to be the directory for [cloudfoundry/bosh-deployment](https://github.com/cloudfoundry/bosh-deployment).
   ```
   # EXAMPLE ONLY!
   $ export BOSH_DEPLOYMENT_DIR=~/Documents/Source/bosh-deployment
   ```

2. Connectivity settings: Libvirt supports a number of connection types, and the `libvirt_cpi` supports 3 of them. Create the variable file  `my-settings.yml` with your specific settings (for ease of copy/paste as well since the `.gitignore` is set for this file name).

   a. `socket`: Connect by Unix socket. The simplest method, but mostly useless for standing up a foundation. If you're just hacking at the CPI code and standing up the Director this is sufficient, and it's quick and easy:
   ```
   libvirt_socket: <socket here, likely /var/run/libvirt/libvirt-sock>
   ```
   Set the `LIBVIRT_CONNECTIVITY` variable:
   ```
   $ export LIBVIRT_CONNECTIVITY=libvirt-socket.yml
   ```
   b. `tcp`: Unsecured access via TCP. See the [Libvirt "Remote support"](https://libvirt.org/remote.html) page for details on how to set up TCP.
   ```
   libvirt_host: <host name or IP address>
   libvirt_port: '<port,likely 16509>'  # must be a string
   ```
   Set the `LIBVIRT_CONNECTIVITY` variable:
   ```
   $ export LIBVIRT_CONNECTIVITY=libvirt-tcp.yml
   ```
   c. `tls`: Secured access via TLS. See the [Libvirt "Remote support"](https://libvirt.org/remote.html) page for details on how to set up TLS connections and generate the various certificates.
   ```
   libvirt_host: <host name or IP address>
   libvirt_port: '<port,likely 16514>'  # must be a string
   libvirt_client_cert: |
     -----BEGIN CERTIFICATE-----
     <snip>
     -----END CERTIFICATE-----
   libvirt_client_key: |
     -----BEGIN RSA PRIVATE KEY-----
     <snip>
     -----END RSA PRIVATE KEY-----
   libvirt_server_ca: |
     -----BEGIN CERTIFICATE-----
     <snip>
     -----END CERTIFICATE-----
   ```
   Set the `LIBVIRT_CONNECTIVITY` variable:
   ```
   $ export LIBVIRT_CONNECTIVITY=libvirt-tls.yml
   ```

3. Create a release. Only needed for development.
   ```
   $ bosh create-release --force --tarball $PWD/cpi
   ```

4. Deploy the BOSH Director.
   ```
   $ bosh create-env ${BOSH_DEPLOYMENT_DIR}/bosh.yml \
       --ops-file=${BOSH_DEPLOYMENT_DIR}/jumpbox-user.yml \
       --ops-file=${BOSH_DEPLOYMENT_DIR}/misc/cpi-resize-disk.yml \
       --ops-file=${BOSH_DEPLOYMENT_DIR}/misc/dns.yml \
       --ops-file=${BOSH_DEPLOYMENT_DIR}/bbr.yml \
       --ops-file=${BOSH_DEPLOYMENT_DIR}/uaa.yml \
       --ops-file=${BOSH_DEPLOYMENT_DIR}/credhub.yml \
       --ops-file=manifests/libvirt-cpi.yml \
       --ops-file=manifests/${LIBVIRT_CONNECTIVITY} \
       --ops-file=manifests/openstack-stemcell.yml \
       --state=state.json \
       --vars-store=bosh-creds.yml \
       --vars-file=my-settings.yml \
       --vars-file=manifests/bosh-vars.yml \
       --vars-file=manifests/libvirt-kvm-vars.yml
   ```
   Notes:
   * For development, add `--ops-file=manifests/libvirt-cpi-dev.yml \` (_after_ the `libvirt-cpi.yml` declaration) to bring in the local copy of the Libvirt CPI.
   * `jumpbox-user.yml` operations file gives access to a user on the BOSH Director that has access to the `root` account. Useful for development work - feel free to leave it off.
   * `cpi-resize-disk.yml` indicates this CPI is able to resize a disk natively. Untested at this time. Feel free to leave it off; note that resizing means BOSH mounts two disks and copies files between the disks, and with current configuration that may not work.
   * `dns.yml` changes the default DNS based on the `internal_dns` entry. Leave it off unless you actually need it.
   * There are two sets of variable files for hypervisor selection: `kvm` and `qemu`, all of which are using the Openstack stemcell. From the Libvirt documentation, `kvm` is likely the best for performance reasons (`openstack-kvm-vars.yml`) as `qemu` virtualizes the entire CPU.
   * With current releases of [bosh-deployment](https://github.com/cloudfoundry/bosh-deployment), the [reboot bug](https://github.com/cloudfoundry/bosh/issues/2131) has been resolved with more current versions of BPM. If you must use an older version of BOSH and the reboot deal is a problem, see the [bosh-reboot-patch](https://github.com/a2geek/bosh-reboot-patch) as it may prove useful. The best option is to get a more current installation of BOSH.
   * Throttling can be enabled. Default is no throttling. Use the `--ops-file=manifests/throttle-noop.yml \` ops file to explicitly enable no-op throttling (default) or `--ops-file=manifests/throttle-file-lock.yml \` to enable file locking. The throttling is achieved by managing the number of _disks_ being created at any time. For a single host, the file I/O appears to be the weak point. Hosts with NVMe or fast SSD's likely will not require throttling. For old style spinning disks, enable the file throttle so save gnashing of teeth, not to mention timeouts!
