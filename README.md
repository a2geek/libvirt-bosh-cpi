# libvirt-bosh-cpi
A Go BOSH CPI for the [libvirt virtualization API](https://libvirt.org/).

## Status

BOSH director can be stood up.

Known defects:
* VMs, when deleted, are only shut down. (Disks are detached and removed.)
* If a VM has a persistent disk attached when it is deleted, that disk also gets deleted. Code patch in place but untested. Likely only an issue when standing up the BOSH director itself(?) since the detach disk method should be called by BOSH itself.

Known TODOs:
* Agent (dynamic) configuration needs to be setup. Working on setting up a configuration disk, currently hardcoded to the [OpenStack settings](https://github.com/cloudfoundry/bosh-linux-stemcell-builder/blob/master/stemcell_builder/stages/bosh_openstack_agent_settings/apply.sh):
  ```
  {
    "Type": "ConfigDrive",
    "DiskPaths": [
      "/dev/disk/by-label/CONFIG-2",
      "/dev/disk/by-label/config-2"
    ],
    "MetaDataPath": "ec2/latest/meta-data.json",
    "UserDataPath": "ec2/latest/user-data"
  },
  ```
* Network is currently assigned via DHCP in the Libvirt settings. Investigate if this can be altered to be configured by the agent.
* Disks are assigned statically; thus more than one of a type will fail. Current scheme:
  * `/dev/vda`: boot disk
  * `/dev/vdb`: ephemeral disk (optional?)
  * `/dev/vdc`: config disk
  * `/dev/vdd`: persistent disk (optional).

## Tinkering

Expect things to not work. This is all setup for Ubuntu Bionic (18.04) using QEMU/KVM for the Libvirt virtualization component.

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
   $ export LIBVIRT_CONNECTIVITY=libvirt_socket.yml
   ```
   b. `tcp`: Unsecured access via TCP. See the [Libvirt "Remote support"](https://libvirt.org/remote.html) page for details on how to set up TCP.
   ```
   libvirt_host: <host name or IP address>
   libvirt_port: '<port,likely 16509>'  # must be a string
   ```
   Set the `LIBVIRT_CONNECTIVITY` variable:
   ```
   $ export LIBVIRT_CONNECTIVITY=libvirt_tcp.yml
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
   $ export LIBVIRT_CONNECTIVITY=libvirt_tls.yml
   ```

3. Create a release. Only needed first time and if there are any code changes.
   ```
   $ bosh create-release --force --tarball $PWD/cpi
   ```

4. Deploy the BOSH Director.
   ```
   $ bosh create-env ${BOSH_DEPLOYMENT_DIR}/bosh.yml \
       --ops-file=${BOSH_DEPLOYMENT_DIR}/jumpbox-user.yml \
       --ops-file=${BOSH_DEPLOYMENT_DIR}/misc/cpi-resize-disk.yml \
       --ops-file=manifests/libvirt_cpi.yml \
       --ops-file=manifests/libvirt_qemu_kvm.yml \
       --ops-file=manifests/${LIBVIRT_CONNECTIVITY} \
       --state=state.json \
       --vars-store=bosh-creds.yml \
       --vars-file=my-settings.yml \
       --vars-file=manifests/bosh-vars.yml
   ```
   Notes:
   * `jumpbox-user.yml` operations file gives access to a user on the BOSH Director that has access to the `root` account. Useful for development work - feel free to leave it off.
   * `cpi-resize-disk.yml` indicates this CPI is able to resize a disk natively. Untested at this time. Feel free to leave it off; note that resizing means BOSH mounts two disks and copies files between the disks, and with current configuration that may not work.

## Setup

Storage space must be allocated to store stemcells and disk images used by the BOSH CPI. 

The following XML defines the bosh pool to be used, assuming that the `default` pool will not be used. See the [libvirt documentation for the storage pool XML](https://libvirt.org/formatstorage.html) for details.

```
<pool type='dir'>
  <name>bosh-pool</name>
  <target>
    <path>/var/lib/libvirt/bosh-images</path>
    <permissions>
      <mode>0711</mode>
      <owner>0</owner>
      <group>0</group>
    </permissions>
  </target>
</pool>
```

Once the location, name, and settings are completed, the following commands will setup the pool and ultimately list details for confirmation:

```
# Run all of these from 'virsh'
pool-define ./libvirt-bosh-storage-pool.xml
pool-build bosh-pool
pool-start bosh-pool
pool-autostart bosh-pool
pool-list
pool-info bosh-pool
```

Sample session:

```
$ virsh
Welcome to virsh, the virtualization interactive terminal.

Type:  'help' for help with commands
       'quit' to quit

virsh # pool-define ./libvirt-bosh-storage-pool.xml
Pool bosh-pool defined from ./libvirt-bosh-storage-pool.xml

virsh # pool-build bosh-pool
Pool bosh-pool built

virsh # pool-start bosh-pool
Pool bosh-pool started

virsh # pool-autostart bosh-pool
Pool bosh-pool marked as autostarted

virsh # pool-list
 Name                 State      Autostart 
-------------------------------------------
 bosh-pool            active     yes       
 default              active     yes       
 Downloads            active     yes       
 tmp                  active     yes       

virsh # pool-info bosh-pool
Name:           bosh-pool
UUID:           0d97e563-5ea4-4aa0-b701-394c01e4acd5
State:          running
Persistent:     yes
Autostart:      yes
Capacity:       437.95 GiB
Allocation:     173.19 GiB
Available:      264.76 GiB

virsh # exit
```

## Developing

At the command-line, from the source directory a compile can be done:

```
$ cd src
$ go build -o ../main main/main.go
```

# References

* [Libvirt "Remote support"](https://libvirt.org/remote.html) 