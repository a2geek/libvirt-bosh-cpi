# Stemcell Compatibility

Development is being done against the Openstack stemcell. If, for some reason, another type of stemcell is needed, support can be expanded.

Current Agent configurations support `ConfigDrive` only, but `CDROM` could be added relatively easily.  `HTTP`, `File`, and `InstanceMetadata` are unlikely to be added.

These are all presumed to be run against some valid Libvirt provider. Thus, there is a chance that the Azure stemcell will work under Libvirt running on KVM.  Hypothetically.

Note that, when trying to identify stemcells, every single one ran under Libvirt/QEMU _except_ for Warden which ships as a GZipped TAR file.  Ultimately, Openstack was settled upon (for no technical reason).  vSphere was hard to confirm due to the initial stemcell password being something different.

## Openstack

```
stemcell:
  formats: [ "openstack-qcow2", "openstack-raw" ]
  type: ConfigDrive
  label: "config-2"
  metadata_path: "ec2/latest/meta-data.json"
  userdata_path: "ec2/latest/user-data"
```

[Source](https://github.com/cloudfoundry/bosh-linux-stemcell-builder/blob/master/stemcell_builder/stages/bosh_openstack_agent_settings/apply.sh)

## Azure

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

## vSphere

> NOTE: When tested with a vSphere stemcell, the CPI failed unable to attach `/dev/sr0`.

```
stemcell:
  formats: [ "vsphere-ova", "vsphere-ovf" ]
  type: CDROM
  filename: env
```

[Source](https://github.com/cloudfoundry/bosh-linux-stemcell-builder/blob/master/stemcell_builder/stages/bosh_vsphere_agent_settings/apply.sh)
