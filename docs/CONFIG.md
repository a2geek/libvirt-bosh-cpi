# Libvirt Configuration

There are a few things to configure with Libvirt with respect to BOSH.

1. Storage. A storage pool is required for the BOSH stemcells and all the VM disks. When tinkering, feel free to use default.
2. Network. The default network is managed by Libvirt and will prevent incoming connections from reaching the VM's. Therefore, an bridge needs to be created manually and a related network created in Libvirt.
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

Instead of using the `default` network, I've created a dedicated network and parcel it out for static IP's as well as dynamic IP's for simple management.

Ubuntu Server now comes with [Netplan](https://netplan.io/) installed. Backup the existing `/etc/netplan/50-cloud-init.yaml` and replace it with something like this:

```
$ cat /etc/netplan/50-cloud-init.yaml
network:
  version: 2
  renderer: networkd

  ethernets:
    enp0s25:
      dhcp4: true

  bridges:
    boshbr0:
      addresses: 
      - 192.168.124.1/24
      interfaces:
      - vlan15

  vlans:
    vlan15:
      accept-ra: no
      id: 15
      link: enp0s25
```

To create the bridge, use `netplan apply` and verify via `brctl show` and `ifconfig` that things look as expected.

In order for network packets to make it into the VMs, you'll need to update the forward rules. In my case, I'm just forwarding to all IP addresses and allowing all VMs to talk amongst themselves. If you need or want to lock these down more, you should be able to lock down ingress by IP or a smaller CIDR block (all VMs should be able to talk amongst themselves).

> For a quick overview of `iptables`, this [stackoverflow answer](https://stackoverflow.com/questions/12945233/iptables-forward-and-input) was helpful for me.

```
# iptables -A FORWARD -d 192.168.124.0/24 -o boshbr0 -j ACCEPT
# iptables -A FORWARD -s 192.168.124.0/24 -i boshbr0 -j ACCEPT
# iptables -A FORWARD -s 192.168.124.0/24 -d 192.168.124.0/24 -i boshbr0 -o boshbr0 -j ACCEPT
# iptables -L FORWARD
Chain FORWARD (policy ACCEPT)
target     prot opt source               destination         
<snip>
ACCEPT     all  --  anywhere             192.168.124.0/24    
ACCEPT     all  --  192.168.124.0/24     anywhere   
ACCEPT     all  --  192.168.124.0/24     192.168.124.0/24    
<snip>
```

> Note: Be certain you get the device names correct. It doesn't simply break in a black-or-white manner if these get fudged up! It may work but not well...

Additionally, BOSH needs to get out to the internet and download software. I found that I needed to masquerade the VM packets (again, I am not a network person, so...):
```
## Do not masquerade to these reserved address blocks.
# iptables -t nat -A POSTROUTING -s 192.168.124.0/24 -d 224.0.0.0/24 -j RETURN
# iptables -t nat -A POSTROUTING -s 192.168.124.0/24 -d 255.255.255.255/32 -j RETURN
## Masquerade all packets going from VMs to the LAN/Internet.
# iptables -t nat -A POSTROUTING -s 192.168.124.0/24 ! -d 192.168.124.0/24 -p tcp -j MASQUERADE --to-ports 1024-65535
# iptables -t nat -A POSTROUTING -s 192.168.124.0/24 ! -d 192.168.124.0/24 -p udp -j MASQUERADE --to-ports 1024-65535
# iptables -t nat -A POSTROUTING -s 192.168.124.0/24 ! -d 192.168.124.0/24 -j MASQUERADE
# iptables -t nat -L
Chain PREROUTING (policy ACCEPT)
target     prot opt source               destination         

Chain INPUT (policy ACCEPT)
target     prot opt source               destination         

Chain OUTPUT (policy ACCEPT)
target     prot opt source               destination         

Chain POSTROUTING (policy ACCEPT)
target     prot opt source               destination         
RETURN     all  --  192.168.124.0/24     base-address.mcast.net/24 
RETURN     all  --  192.168.124.0/24     255.255.255.255     
MASQUERADE  tcp  --  192.168.124.0/24    !192.168.124.0/24     masq ports: 1024-65535
MASQUERADE  udp  --  192.168.124.0/24    !192.168.124.0/24     masq ports: 1024-65535
MASQUERADE  all  --  192.168.124.0/24    !192.168.124.0/24    
```
(Pulled from [here](https://jamielinux.com/docs/libvirt-networking-handbook/custom-nat-based-network.html).)

Once the `iptables` configuration is complete/correct, it can be persisted. See [this discussion](https://askubuntu.com/questions/66890/how-can-i-make-a-specific-set-of-iptables-rules-permanent) for more details.

```
# cat /etc/iptables.conf 
# Generated by iptables-save v1.6.1 on Sun Oct 13 22:17:37 2019
*nat
:PREROUTING ACCEPT [56:10387]
:INPUT ACCEPT [0:0]
:OUTPUT ACCEPT [20:1656]
:POSTROUTING ACCEPT [20:1656]
-A POSTROUTING -s 192.168.124.0/24 -d 224.0.0.0/24 -j RETURN
-A POSTROUTING -s 192.168.124.0/24 -d 255.255.255.255/32 -j RETURN
-A POSTROUTING -s 192.168.124.0/24 ! -d 192.168.124.0/24 -p tcp -j MASQUERADE --to-ports 1024-65535
-A POSTROUTING -s 192.168.124.0/24 ! -d 192.168.124.0/24 -p udp -j MASQUERADE --to-ports 1024-65535
-A POSTROUTING -s 192.168.124.0/24 ! -d 192.168.124.0/24 -j MASQUERADE
COMMIT
# Completed on Sun Oct 13 22:17:37 2019
# Generated by iptables-save v1.6.1 on Sun Oct 13 22:17:37 2019
*filter
:INPUT ACCEPT [0:0]
:FORWARD ACCEPT [0:0]
:OUTPUT ACCEPT [828:146927]
-A INPUT -i lo -j ACCEPT
-A INPUT -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
-A INPUT -p tcp -m tcp --dport 16514 -j ACCEPT
-A INPUT -p tcp -m tcp --dport 22 -j ACCEPT
-A INPUT -j DROP
-A FORWARD -d 192.168.124.0/24 -o boshbr0 -j ACCEPT
-A FORWARD -s 192.168.124.0/24 -i boshbr0 -j ACCEPT
-A FORWARD -s 192.168.124.0/24 -d 192.168.124.0/24 -i boshbr0 -o boshbr0 -j ACCEPT
-A FORWARD -j DROP
COMMIT
# Completed on Sun Oct 13 22:17:37 2019
```

Update `/etc/rc.local` with this to reload:
```
#!/bin/bash

# Load iptables rules from this file
iptables-restore < /etc/iptables.conf
```
(Note: If using `systemd` - Ubuntu does - `/etc/rc.local` must start with a [shebang](https://en.wikipedia.org/wiki/Shebang_(Unix)).)

Finally, create the bridge in Libvirt:

```
$ cat network.xml 
<network connections='25'>
  <name>bosh</name>
  <forward mode='bridge'/>
  <bridge name='boshbr0' />
</network>
$ virsh net-define network.xml 
Network bosh defined from network.xml

$ virsh net-start bosh
Network bosh started

$ virsh net-autostart bosh
Network bosh marked as autostarted

```

## Remote access

Follow the documentation provided on the Libvirt site to generate certificates:
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

### Ubuntu 20.04

To enable TLS, likely all that needs to be setup is:

```
# systemctl start libvirtd-tls.socket
```

### Ubuntu 18.04

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
