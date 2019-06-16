package cpi

// LibvirtVMCloudProps represents the VMCloudProps supplied in the BOSH Cloud properties
//  hash specified in the deployment manifest under the VM's resource pool.
type LibvirtVMCloudProps struct {
	// Number of (virtual) CPUs
	CPU uint `json:"cpu"`
	// Memory sized in megabytes
	Memory uint `json:"memory"`
	// EphemeralDisk sized in megabytes.
	EphemeralDisk uint64 `json:"ephemeral_disk"`
}

// LibvirtStemcellCloudProps represents the cloud properties supplied with stemcells.
type LibvirtStemcellCloudProps struct {
	Name           string `json:"name"`
	Version        string `json:"version"`
	Infrastructure string `json:"infrastructure"`
	Hypervisor     string `json:"hypervisor"`
	// Size of the stemcell image in MiB's.
	Disk            uint64 `json:"disk"`
	DiskFormat      string `json:"disk_format"`
	ContainerFormat string `json:"container_format"`
	OsType          string `json:"os_type"`
	OsDistro        string `json:"os_distro"`
	Architecture    string `json:"architecture"`
	AutoDiskConfig  bool   `json:"auto_disk_config"`
}
