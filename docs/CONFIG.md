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

Follow the documenation provided on the Libvirt site:
* [Libvirt "Remote support"](https://libvirt.org/remote.html)
