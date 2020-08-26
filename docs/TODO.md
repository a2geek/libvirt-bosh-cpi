# TODOs

* Network is currently assigned via DHCP in the Libvirt settings. Investigate if this can be altered to be configured by the agent.
* Network only supports `manual` networks at this time. (Have not run across any other type at this point.)
* Can the stemcell be interrogated for the stemcell configuration settings? currently is hardcoded. Maybe the metdata section can be harnessed to persist _with_ the boot disk or stemcell. Is this worth it?
* Restructure "manager" api.
