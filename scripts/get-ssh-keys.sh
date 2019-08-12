#!/bin/bash

[ -f vcap-private-key.pem ] && chmod 660 vcap-private-key.pem
[ -f jumpbox-private-key.pem ] && chmod 660 jumpbox-private-key.pem

bosh int bosh-creds.yml --path /vm_ssh_key/private_key > vcap-private-key.pem
bosh int bosh-creds.yml --path /jumpbox_ssh/private_key > jumpbox-private-key.pem

chmod 400 vcap-private-key.pem jumpbox-private-key.pem
