# Deployments

Here are some experiments for deployments. The default cloud config should be set for them. If you're going to run these "for real", review the actual deployment instructions - especially Cloud Foundry, as it's been reduced to fit on a smaller development machine.

Be certain the BOSH environment has been setup in the current shell:
```
$ source scripts/bosh-env.sh
```

Make sure the cloud config exists (or is current):
```
$ bosh -n update-cloud-config manifests/cloud-config.yml
```

Upload a stemcell (check for current Openstack stemcells at [bosh.io](https://bosh.io/stemcells/bosh-openstack-kvm-ubuntu-xenial-go_agent)):
```
$ bosh upload-stemcell --sha1 9c6153c5a41b48e5833b1f25fced4e06fa6d6ba1 \
    https://bosh.io/d/stemcells/bosh-openstack-kvm-ubuntu-xenial-go_agent?v=456.51
```

# Samples

These deployments can be used as a starting point!

* [Zookeeper](ZOOKEEPER.md)
* [Postgres](POSTGRES.md)
* [Concourse](CONCOURSE.md)
* [Cloud Foundry](CF.md)
* [CF MySQL](CF_MYSQL.md) - brokered service
* [CF Redis](CF_REDIS.md) - brokered service
* [Blacksmith](BLACKSMITH.md) - multiple brokered services, including clustered services
* [Kubo/Kubernetes](KUBO.md)
