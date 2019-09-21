# MySQL

The MySQL database proxies require some form of load balancer. At this point, there is none pre-provisioned in Libvirt (or KVM/QEMU), so something needs to be created.

This command stands up an NGINX proxy. Note that the IP addresses in the file were pulled from the deployment, so those may need to be updated. 

This configuration does not include any form of health check as that appears to be part of NGINX Plus, their commercial offering. :-(

```
$ bosh -n -d mysql-nginx-proxy deploy manifests/mysql-nginx-proxy.yml \
    --vars-store=mysql-proxy-creds.yml
```

The `mysql-vars.yml` file has a `cf_mysql_host` (default of `mysql.sys.mycf.lan`) variable pointing to the DNS entry for the MySQL database proxy (currently NGINX). Be certain to setup your DNS (or `/etc/hosts`) file appropriately.

Set `CF_MYSQL_DEPLOYMENT` to a local copy of the [CF MySQL Deployment](https://github.com/cloudfoundry/cf-mysql-deployment/) directory.
```
$ export CF_MYSQL_DEPLOYMENT=~/Documents/Source/cf-mysql-deployment/
```

Ensure you are on the `master` branch (it defaults to `develop`).
```
$ pushd ${CF_MYSQL_DEPLOYMENT}
$   git checkout master
$ popd
```

Get the CredHub admin credentials and login to CredHub. (Only needed if you aren't already logged in, obviously.)
```
$ source scripts/credhub-env.sh
$ credhub login
Setting the target url: https://192.168.123.7:8844
Login Successful
```

Add the CF `admin` password into CredHub to make it available.
```
$ bosh int cf-creds.yml --path /cf_admin_password | \
    credhub set -n /libvirt/cf-mysql/cf_admin_password -t password
```

Deploy!
```
$ bosh -n -d cf-mysql deploy ${CF_MYSQL_DEPLOYMENT}/cf-mysql-deployment.yml \
    -o ${CF_MYSQL_DEPLOYMENT}/operations/add-broker.yml \
    -o ${CF_MYSQL_DEPLOYMENT}/operations/register-proxy-route.yml \
    -o ${CF_MYSQL_DEPLOYMENT}/operations/xenial-stemcell.yml \
    -o ${CF_MYSQL_DEPLOYMENT}/operations/configure-broker-load-balancer.yml \
    -l manifests/mysql-vars.yml \
    --vars-store=mysql-creds.yml
```

Run the broker registration errand:
```
$ bosh -d cf-mysql run-errand broker-registrar
```

Once an org and space have been selected, a database can be created:
```
$ cf create-service p-mysql 10mb test-db
Creating service instance test-db in org rob / space dev as admin...
OK

$ cf service test-db
Showing info of service test-db in org rob / space dev as admin...

name:             test-db
service:          p-mysql
tags:             
plan:             10mb
description:      MySQL databases on demand
documentation:    https://github.com/cloudfoundry/cf-mysql-release/blob/master/README.md
dashboard:        https://cf-mysql.sys.mycf.lan/manage/instances/a6d60f57-523e-4cd1-8f28-498beaa726ad
service broker:   p-mysql

Showing status of last operation from service test-db...

status:    create succeeded
message:   
started:   2019-09-21T19:57:43Z
updated:   2019-09-21T19:57:43Z

There are no bound apps for this service.

Upgrades are not supported by this broker.
```

The dashboard also should work with an appropriately provisioned id.
