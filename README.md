# libvirt-bosh-cpi
A Go BOSH CPI for the [libvirt virtualization API](https://libvirt.org/).

## Status

BOSH director can be stood up. Stemcell can now be uploaded.

Known defects:
* If a VM has a persistent disk attached when it is deleted, that disk also gets deleted. Likely only an issue when standing up the BOSH director itself(?) since the detach disk method should be called by BOSH itself.

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
       --ops-file=${BOSH_DEPLOYMENT_DIR}/misc/dns.yml \
       --ops-file=manifests/libvirt_cpi.yml \
       --ops-file=manifests/${LIBVIRT_CONNECTIVITY} \
       --state=state.json \
       --vars-store=bosh-creds.yml \
       --vars-file=my-settings.yml \
       --vars-file=manifests/bosh-vars.yml \
       --vars-file=manifests/openstack-kvm-vars.yml
   ```
   Notes:
   * `jumpbox-user.yml` operations file gives access to a user on the BOSH Director that has access to the `root` account. Useful for development work - feel free to leave it off.
   * `cpi-resize-disk.yml` indicates this CPI is able to resize a disk natively. Untested at this time. Feel free to leave it off; note that resizing means BOSH mounts two disks and copies files between the disks, and with current configuration that may not work.
   * `dns.yml` changes the default DNS based on the `internal_dns` entry. Leave it off unless you actually need it.
   * There are two sets of variable files for hypervisor selection: `kvm` and `qemu`, all of which are using the Openstack stemcell. From the Libvirt documentation, `kvm` is likely the best for performance reasons (`openstack-kvm-vars.yml`) as `qemu` virtualizes the entire CPU.

## Deployments

Here are some experiments for deployments. The default cloud config should be set for them.

Be certain the BOSH environment has been setup in the current shell:
```
$ source scripts/bosh-env.sh
```

Make sure the cloud config exists (or is current):
```
$ bosh -n update-cloud-config manifests/cloud-config.yml
```

Upload a stemcell (check for current Openstack stemcells at [bosh.io](https://bosh.io/stemcells/bosh-openstack-kvm-ubuntu-xenial-go_agent)):
```
$ bosh upload-stemcell --sha1 b5f9671591b22602b982fbf4f2320fe971718f7e  https://bosh.io/d/stemcells/bosh-openstack-kvm-ubuntu-xenial-go_agent?v=456.3
```

### Zookeeper

Set `ZOOKEEPER_DIR` to a local copy of the [Zookeeper release](https://github.com/cppforlife/zookeeper-release) directory.

Deploy!
```
$ bosh -n -d zookeeper deploy $ZOOKEEPER_DIR/manifests/zookeeper.yml \
    --vars-store=zookeeper-creds.yml
```

### Postgres

Set `POSTGRES_DIR` to a local copy of the [Postgres release](https://github.com/cloudfoundry/postgres-release) directory.

Upload a Postgres release...
```
$ bosh upload-release https://bosh.io/d/github.com/cloudfoundry/postgres-release
```

Deploy!
```
$ bosh -n -d postgres deploy $POSTGRES_DIR/templates/postgres.yml \
    -o $POSTGRES_DIR/templates/operations/add_static_ips.yml \
    -o $POSTGRES_DIR/templates/operations/set_properties.yml \
    -o $POSTGRES_DIR/templates/operations/use_bbr.yml \
    --vars-store=postgres.yml \
    -l manifests/postgres-vars.yml
```

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

## Utilities

There are a few scripts in the `scripts` folder.

1. BOSH environment variables for the current Director can be setup with `source scripts/bosh-env.sh`.
2. The SSH keys can be extracted with `scripts/get-ssh-keys.sh`. These will create two PEM files: `vcap-private-key.pem` and `jumpbox-private-key.pem`. Usage is the usual SSH mechanisms like `ssh -i jumpbox-private-key.pem jumpbox@your-director`. If you change the director, any trusts need to be resolved as usual.

## Developing

At the command-line, from the source directory a compile can be done:

```
$ cd src
$ go build -o ../main main/main.go
```

## References

* [Libvirt "Remote support"](https://libvirt.org/remote.html)
* [BOSH Agent MetadataContentsType](https://godoc.org/github.com/cloudfoundry/bosh-agent/infrastructure#MetadataContentsType)
* [BOSH Agent UserDataContentsType](https://godoc.org/github.com/cloudfoundry/bosh-agent/infrastructure#UserDataContentsType)
