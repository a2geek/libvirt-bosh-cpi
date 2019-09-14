#!/bin/bash

pushd $(dirname ${BASH_SOURCE})/.. > /dev/null
    export CF_USERNAME=admin
    export CF_PASSWORD=$(bosh int cf-creds.yml --path /cf_admin_password)
popd > /dev/null
