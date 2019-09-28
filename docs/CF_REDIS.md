# Redis

Set `CF_REDIS_DEPLOYMENT` to a local copy of the [Shared Redis Release](https://github.com/pivotal-cf/shared-redis-release/) directory.
```
$ export CF_REDIS_DEPLOYMENT=~/Documents/Source/shared-redis-release/
```

Note that there are a number of releases documented as being required, please see [[Shared Redis Release](https://github.com/pivotal-cf/shared-redis-release/) if additional need to be uploaded.
```
$ bosh upload-release http://bosh.io/d/github.com/cloudfoundry/syslog-release
```

There is some pre-work that needs to be done since the release isn't published anywhere:
```
$ pushd $CF_REDIS_DEPLOYMENT
$   git submodule update --init --recursive
$   bosh create-release
$   bosh upload-release
$ popd
```

Add the CF `admin` password into CredHub to make it available.
```
$ bosh int cf-creds.yml --path /cf_admin_password | \
    credhub set -n /libvirt/cf-redis/cf_password -t password
```

Generate and store the redis `broker` password in CredHub:
```
$ credhub generate -n /libvirt/cf-redis/broker_password -t password
```

Deploy!
```
$ bosh -n -d cf-redis deploy ${CF_REDIS_DEPLOYMENT}/manifest/deployment.yml \
    -l manifests/redis-vars.yml \
    --vars-store=redis-creds.yml
```

And, finally register that broker:
```
$ bosh -d cf-redis run-errand broker-registrar
```

Success!
```
$ cf marketplace
Getting services from marketplace in org rob / space dev as admin...
OK

service   plans        description                                  broker
p-redis   shared-vm    Redis service to provide a key-value store   cf-redis-broker

TIP: Use 'cf marketplace -s SERVICE' to view descriptions of individual plans of a given service.
```

Note that the current configuration does not included dedicated nodes, but they can be added by tweaking the variables.
