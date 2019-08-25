# libvirt-bosh-cpi
A Go BOSH CPI for the [libvirt virtualization API](https://libvirt.org/).

## Status

Many things do work at this point. Most likely a BOSH deployment will successfully _deploy_. Longer term management (resizing disks or snapshotting vms, for instance) either are untested or unimplemented at this point.

Known TODOs:
* Network is currently assigned via DHCP in the Libvirt settings. Investigate if this can be altered to be configured by the agent.
* Disks are assigned statically; thus more than one of a type will fail. Current scheme:
  * `/dev/vda`: boot disk
  * `/dev/vdb`: ephemeral disk (optional?)
  * `/dev/vdc`: config disk
  * `/dev/vdd`: persistent disk (optional).
* Restructure this README as it's gotten both overwhelming and uninformative. Yes, that's a thing. ;-)

## Stemcell Compatibility

Development is being done against the Openstack stemcell. If, for some reason, another type of stemcell is needed, support can be expanded.

Current Agent configurations support `ConfigDrive` only, but `CDROM` could be added relatively easily.  `HTTP`, `File`, and `InstanceMetadata` are unlikely to be added.

These are all presumed to be run against some valid Libvirt provider. Thus, there is a chance that the Azure stemcell will work under Libvirt running on KVM.  Hypothetically.

Note that, when trying to identify stemcells, every single one ran under Libvirt/QEMU _except_ for Warden which ships as a GZipped TAR file.  Ultimately, Openstack was settled upon (for no technical reason).  vSphere was hard to confirm due to the initial stemcell password being something different.

### Openstack

```
stemcell:
  formats: [ "openstack-qcow2", "openstack-raw" ]
  type: ConfigDrive
  label: "config-2"
  metadata_path: "ec2/latest/meta-data.json"
  userdata_path: "ec2/latest/user-data"
```

[Source](https://github.com/cloudfoundry/bosh-linux-stemcell-builder/blob/master/stemcell_builder/stages/bosh_openstack_agent_settings/apply.sh)

### Azure

> NOTE: The `label` is too long for current config disk capabilities. At a minimum, there may be other issues.

```
stemcell:
  formats: [ "azure-vhd" ]
  type: ConfigDrive
  label: azure_cfg_dsk
  metadata_path: "configs/MetaData"
  userdata_path: "configs/UserData"
```

[Source](https://github.com/cloudfoundry/bosh-linux-stemcell-builder/blob/master/stemcell_builder/stages/bosh_azure_agent_settings/apply.sh)

### vSphere

> NOTE: When tested with a vSphere stemcell, the resulting VM did not respond. 

```
stemcell:
  formats: [ "vsphere-ova", "vsphere-ovf" ]
  type: CDROM
  filename: env
```

[Source](https://github.com/cloudfoundry/bosh-linux-stemcell-builder/blob/master/stemcell_builder/stages/bosh_vsphere_agent_settings/apply.sh)

## Tinkering

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
   * `jumpbox-user.yml` operations file gives access to a user on the BOSH Director that has access to the `root` account. Useful for development work - feel free to leave it off.
   * `cpi-resize-disk.yml` indicates this CPI is able to resize a disk natively. Untested at this time. Feel free to leave it off; note that resizing means BOSH mounts two disks and copies files between the disks, and with current configuration that may not work.
   * `dns.yml` changes the default DNS based on the `internal_dns` entry. Leave it off unless you actually need it.
   * There are two sets of variable files for hypervisor selection: `kvm` and `qemu`, all of which are using the Openstack stemcell. From the Libvirt documentation, `kvm` is likely the best for performance reasons (`openstack-kvm-vars.yml`) as `qemu` virtualizes the entire CPU.
   * With current releases of [bosh-deployment](https://github.com/cloudfoundry/bosh-deployment), the [reboot bug](https://github.com/cloudfoundry/bosh/issues/2131) has been resolved with more current versions of BPM. If you must use an older version of BOSH and the reboot deal is a problem, see the [bosh-reboot-patch](https://github.com/a2geek/bosh-reboot-patch) as it may prove useful. The best option is to get a more current installation of BOSH.

## Deployments

Here are some experiments for deployments. The default cloud config should be set for them. If you're going to run these "for real", review the actual deployment instructions - especially Cloud Foundry, as it's been reduced to fit on a smaller development machine.

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
    --vars-store=postgres-creds.yml \
    -l manifests/postgres-vars.yml
```

### Concourse

Set `CONCOURSE_DIR` to a local copy of the [Concourse BOSH deployment](https://github.com/concourse/concourse-bosh-deployment) directory.

Deploy!
```
$ bosh -n -d concourse deploy $CONCOURSE_DIR/cluster/concourse.yml \
    -o $CONCOURSE_DIR/cluster/operations/basic-auth.yml \
    -o $CONCOURSE_DIR/cluster/operations/static-web.yml \
    -o $CONCOURSE_DIR/cluster/operations/privileged-http.yml \
    -l $CONCOURSE_DIR/versions.yml \
    --vars-store=concourse-creds.yml \
    -l manifests/concourse-vars.yml
```

Concourse will be available at http://192.168.123.250 (assuming all the network stuff fits with your setup).

### Cloud Foundry

Set `CF_DEPLOYMENT_DIR` to a local copy of the [Cloud Foundry deployment](https://github.com/cloudfoundry/cf-deployment/) directory.

Add the DNS runtime config:
```
$ bosh update-runtime-config $BOSH_DEPLOYMENT_DIR/runtime-configs/dns.yml --name dns
```

Deploy!
```
$ bosh -n -d cf deploy $CF_DEPLOYMENT_DIR/cf-deployment.yml \
    -o $CF_DEPLOYMENT_DIR/operations/scale-to-one-az.yml \
    -o $CF_DEPLOYMENT_DIR/operations/set-router-static-ips.yml \
    -o $CF_DEPLOYMENT_DIR/operations/use-compiled-releases.yml \
    -o $CF_DEPLOYMENT_DIR/operations/use-latest-stemcell.yml \
    -l manifests/cloudfoundry-vars.yml \
    --vars-store=cf-creds.yml
```

Note that there is a requirement for DNS resolution to `*.sys.mypcf.lan` as currently configured. `/etc/hosts` can be used as a hack for validation.

```
$ cat /etc/hosts | grep mypcf
192.168.123.252 api.sys.mypcf.lan login.sys.mypcf.lan sample1.sys.mypcf.lan
```

To get the admin credentials setup:
```
$ export CREDHUB_CLIENT=credhub-admin
$ export CREDHUB_SECRET=$(bosh interpolate ./cf-creds.yml --path=/credhub_admin_client_secret)
$ export CREDHUB_CA_CERT="$(bosh interpolate ./cf-creds.yml --path=/credhub_tls/ca )"$'\n'"$( bosh interpolate ./cf-creds.yml --path=/uaa_ssl/ca)"
```

To login with those credentials:
```
$ cf api --skip-ssl-validation https://api.sys.mypcf.lan
Setting api endpoint to https://api.sys.mypcf.lan...
OK

api endpoint:   https://api.sys.mypcf.lan
api version:    2.139.0
$ cf login 
API endpoint: https://api.sys.mypcf.lan

Email> admin

Password> (paste in cf_admin_password from cf-creds.yml file)
Authenticating...
OK

Targeted org system

API endpoint:   https://api.sys.mypcf.lan (API version: 2.139.0)
User:           admin
Org:            system
Space:          No space targeted, use 'cf target -s SPACE'
```

Finally, create a place to deploy applications:
```
$ cf create-org robstuff
Creating org robstuff as admin...
OK

Assigning role OrgManager to user admin in org robstuff ...
OK

TIP: Use 'cf target -o "robstuff"' to target new org

$ cf target -o robstuff
api endpoint:   https://api.sys.mypcf.lan
api version:    2.139.0
user:           admin
org:            robstuff
No space targeted, use 'cf target -s SPACE'

$ cf create-space np
Creating space np in org robstuff as admin...
OK
Assigning role RoleSpaceManager to user admin in org robstuff / space np as admin...
OK
Assigning role RoleSpaceDeveloper to user admin in org robstuff / space np as admin...
OK

TIP: Use 'cf target -o "robstuff" -s "np"' to target new space

$ cf target -s np
api endpoint:   https://api.sys.mypcf.lan
api version:    2.139.0
user:           admin
org:            robstuff
space:          np
```

Note that `sample1.sys.mypcf.lan` is just for a quick test deploy like this:
```
$ mkdir staticfile-sample
$ cd staticfile-sample
staticfile-sample$ touch Staticfile
staticfile-sample$ cat > index.html
Hello World!
^D
staticfile-sample$ cf push -b staticfile_buildpack -m 32M -p . sample1
staticfile-sample$ cf apps
Getting apps in org robstuff / space np as admin...
OK

name      requested state   instances   memory   disk   urls
sample1   started           1/1         32M      1G     sample1.sys.mypcf.lan
```

## Libvirt Configuration

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
* [Cloud Foundry Deployment Guide](https://github.com/cloudfoundry/cf-deployment/blob/master/texts/deployment-guide.md)
* [Cloud Foundry Cloud Configs](https://github.com/cloudfoundry/cf-deployment/blob/master/texts/on-cloud-configs.md)