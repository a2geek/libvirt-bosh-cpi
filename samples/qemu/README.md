# Experiments

Using `virsh` to test/understand how various libvirt operations work and functional XML.

## 01 Create stemcell

```
$ virsh vol-create --pool default --file 01-create-stemcell.xml
Vol bosh-openstack-kvm-ubuntu-xenial-go_agent-170.14 created from 01-create-stemcell.xml

$ virsh vol-list --pool default
 Name                 Path                                    
------------------------------------------------------------------------------
 bosh-openstack-kvm-ubuntu-xenial-go_agent-170.14 /var/lib/libvirt/images/bosh-openstack-kvm-ubuntu-xenial-go_agent-170.14

$ ls -l /var/lib/libvirt/images/bosh-openstack-kvm-ubuntu-xenial-go_agent-170.14
-rw------- 1 root root 863109120 Jun  7 17:15 /var/lib/libvirt/images/bosh-openstack-kvm-ubuntu-xenial-go_agent-170.14
$ virsh vol-upload --pool default --vol bosh-openstack-kvm-ubuntu-xenial-go_agent-170.14 --file ~/Downloads/stemcells/tmp/root.img

$ virsh vol-info --pool default --vol bosh-openstack-kvm-ubuntu-xenial-go_agent-170.14
Name:           bosh-openstack-kvm-ubuntu-xenial-go_agent-170.14
Type:           file
Capacity:       3.00 GiB
Allocation:     823.12 MiB

```

## 02 Clone root volume from stemcell

```
$ virsh vol-clone --pool default --vol bosh-openstack-kvm-ubuntu-xenial-go_agent-170.14 --newname test-root-volume
Vol test-root-volume cloned from bosh-openstack-kvm-ubuntu-xenial-go_agent-170.14

$ virsh vol-list --pool default
 Name                 Path                                    
------------------------------------------------------------------------------
 bosh-openstack-kvm-ubuntu-xenial-go_agent-170.14 /var/lib/libvirt/images/bosh-openstack-kvm-ubuntu-xenial-go_agent-170.14
 test-root-volume     /var/lib/libvirt/images/test-root-volume

$ virsh vol-info --pool default --vol test-root-volume
Name:           test-root-volume
Type:           file
Capacity:       3.00 GiB
Allocation:     1.69 GiB

```

## 03 Create ephemeral disk

```
$ virsh vol-create --pool default --file 03-create-ephemeral.xml
Vol test-ephemeral-volume created from 03-create-ephemeral.xml

$ virsh vol-list --pool default
 Name                 Path                                    
------------------------------------------------------------------------------
 bosh-openstack-kvm-ubuntu-xenial-go_agent-170.14 /var/lib/libvirt/images/bosh-openstack-kvm-ubuntu-xenial-go_agent-170.14
 test-ephemeral-volume /var/lib/libvirt/images/test-ephemeral-volume
 test-root-volume     /var/lib/libvirt/images/test-root-volume

$ virsh vol-info --pool default --vol test-ephemeral-volume
Name:           test-ephemeral-volume
Type:           file
Capacity:       16.00 GiB
Allocation:     0.00 B

```

## 04 Create VM

Note that the device attachments had to be file/path based and not based on the storage volume or storage volume name.

Also note that the Domain/VM only has a console, meaning [Virt Manager](https://virt-manager.org/) does not display a graphical console.

```
$ virsh define --file 04a-create-vm.xml 
Domain vm-8f7b3bbd-f777-466d-940b-312ecc3c6db6 defined from 04a-create-vm.xml

$ virsh attach-device --domain vm-8f7b3bbd-f777-466d-940b-312ecc3c6db6 --file 04b-add-root-device.xml --current
Device attached successfully

$ virsh attach-device --domain vm-8f7b3bbd-f777-466d-940b-312ecc3c6db6 --file 04c-add-ephemeral-device.xml --current
Device attached successfully

$ virsh attach-device --domain vm-8f7b3bbd-f777-466d-940b-312ecc3c6db6 --file 04d-add-network-device.xml --current
Device attached successfully

$ virsh console vm-8f7b3bbd-f777-466d-940b-312ecc3c6db6
Connected to domain vm-8f7b3bbd-f777-466d-940b-312ecc3c6db6
Escape character is ^]
[    2.600222] ACPI: PCI Interrupt Link [LNKB] enabled at IRQ 10
[    2.834706] ACPI: PCI Interrupt Link [LNKC] enabled at IRQ 11
[    3.293596] ACPI: PCI Interrupt Link [LNKA] enabled at IRQ 10
[    3.297472] Serial: 8250/16550 driver, 32 ports, IRQ sharing enabled
[    3.319360] 00:04: ttyS0 at I/O 0x3f8 (irq = 4, base_baud = 115200) is a 16550A
[    3.340765] Linux agpgart interface v0.103
[    3.357154] loop: module loaded
[    3.364558] scsi host0: ata_piix
[    3.366479] scsi host1: ata_piix
[    3.366999] ata1: PATA max MWDMA2 cmd 0x1f0 ctl 0x3f6 bmdma 0xc0e0 irq 14
<snip>
```
