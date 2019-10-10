# Libvirt Configuration

There are a few things to configure with Libvirt with respect to BOSH.

1. Storage. A storage pool is required for the BOSH stemcells and all the VM disks. When tinkering, feel free to use default.
2. Network. So far, the default network has been sufficient - a class "C" network (~254 IP addresses). Note that this appears to be random, so likely the IP addresses need to be adjusted.
3. Remote access. Ultimately, the BOSH VM will be running within a Libvirt managed VM and will *not* have access to the default Unix socket.

## Storage

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

## Network

Sample of a default network that was generated:

```
virsh # net-list
 Name                 State      Autostart     Persistent
----------------------------------------------------------
 default              active     yes           yes

virsh # net-dumpxml default
<network connections='25'>
  <name>default</name>
  <uuid>bcb85dc8-ad27-4861-9f85-692167fd79fa</uuid>
  <forward mode='nat'>
    <nat>
      <port start='1024' end='65535'/>
    </nat>
  </forward>
  <bridge name='virbr0' stp='on' delay='0'/>
  <mac address='52:54:00:82:d1:6e'/>
  <ip address='192.168.123.1' netmask='255.255.255.0'>
    <dhcp>
      <range start='192.168.123.2' end='192.168.123.254'/>
    </dhcp>
  </ip>
</network>
```

## Remote access

Follow the documenation provided on the Libvirt site to generate certificates:
* [Libvirt "Remote support"](https://libvirt.org/remote.html)

Note that the `serverkey.pem` file likely needs to be in `/etc/pki/libvirt/private/` instead based on Ubuntu configurations, rather that the documented `/etc/pki/libvirt` directory.

Copying files into place:

```
$ sudo mkdir -p /etc/pki/CA /etc/pki/libvirt/private
$ sudo cp cacert.pem /etc/pki/CA/cacert.pem
$ sudo cp servercert.pem /etc/pki/libvirt/
$ sudo cp serverkey.pem /etc/pki/libvirt/private/
```

For permissions, the Internet indicates the directories should be `700` and the files should be `600`:

```
$ sudo chmod 700 /etc/pki /etc/pki/CA /etc/pki/libvirt /etc/pki/libvirt/private
$ sudo find /etc/pki -name "*.pem" -exec chmod 600 {} \;
$ sudo find /etc/pki -ls
 14681727      4 drwx------   4 root     root         4096 Oct 10 02:45 /etc/pki
 14681728      4 drwx------   2 root     root         4096 Oct 10 02:44 /etc/pki/CA
 14681729      4 -rw-------   1 root     root         1432 Oct 10 02:44 /etc/pki/CA/cacert.pem
 14681730      4 drwx------   3 root     root         4096 Oct 10 03:47 /etc/pki/libvirt
 14681731      4 -rw-------   1 root     root         1598 Oct 10 02:45 /etc/pki/libvirt/servercert.pem
 14681635      4 drwx------   2 root     root         4096 Oct 10 03:45 /etc/pki/libvirt/private
 14681643      8 -rw-------   1 root     root         8170 Oct 10 03:45 /etc/pki/libvirt/private/serverkey.pem
```

To enable TLS, the config in `/etc/default/libvirtd` needs to be altered (note `libvirt_opts`):

```
# cat /etc/default/libvirtd
# Defaults for libvirtd initscript (/etc/init.d/libvirtd)
# This is a POSIX shell fragment

# Start libvirtd to handle qemu/kvm:
start_libvirtd="yes"

# options passed to libvirtd, add "-l" to listen on tcp
libvirtd_opts="-l"    # <== uncomment this line

# pass in location of kerberos keytab
#export KRB5_KTNAME=/etc/libvirt/libvirt.keytab

# Whether to mount a systemd like cgroup layout (only
# useful when not running systemd)
#mount_cgroups=yes
# Which cgroups to mount
#cgroups="memory devices"
```

By default, TLS should be enabled and this should place the `*.pem` files into the correct locations.
