# TODOs

* Network is currently assigned via DHCP in the Libvirt settings. Investigate if this can be altered to be configured by the agent.
* Disks are assigned statically; thus more than one of a type will fail. Current scheme:
  * `/dev/vda`: boot disk
  * `/dev/vdb`: ephemeral disk (optional?)
  * `/dev/vdc`: config disk
  * `/dev/vdd`: persistent disk (optional).
* Can the stemcell be interrogated for the stemcell configuration settings? currently is hardcoded. Maybe the metdata section can be harnessed to persist _with_ the boot disk or stemcell. Is this worth it?
* Restructure "manager" api.
