# TODOs

* BUG: Something with disks -- maybe a connection isn't being closed? A reboot of the BOSH VM resolved the error.
  ```
    Task 635 | 19:31:37 | Updating instance web: web/b961839a-660c-419a-9140-4595eb8d9d8c (0) (canary) (00:00:19)
                        L Error: CPI error 'Bosh::Clouds::NotImplemented' with message 'Must call implemented method: failed to dial libvirt: dial tcp: i/o timeout' in 'delete_vm' CPI method (CPI request ID: 'cpi-657334')
  ```
* Network is currently assigned via DHCP in the Libvirt settings. Investigate if this can be altered to be configured by the agent.
* Network only supports `manual` networks at this time. (Have not run across any other type at this point.)
* Disks are assigned statically; thus more than one of a type will fail. Current scheme:
  * `/dev/vda`: boot disk
  * `/dev/vdb`: ephemeral disk (optional?)
  * `/dev/vdc`: config disk
  * `/dev/vdd`: persistent disk (optional).
* Can the stemcell be interrogated for the stemcell configuration settings? currently is hardcoded. Maybe the metdata section can be harnessed to persist _with_ the boot disk or stemcell. Is this worth it?
* Restructure "manager" api.
