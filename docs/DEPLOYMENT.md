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
$ bosh upload-stemcell --sha1 b5f9671591b22602b982fbf4f2320fe971718f7e  https://bosh.io/d/stemcells/bosh-openstack-kvm-ubuntu-xenial-go_agent?v=456.3
```

# Samples

These deployments can be used as a starting point!

* [Zookeeper](ZOOKEEPER.md)
* [Postgres](POSTGRES.md)
* [Concourse](CONCOURSE.md)
* [Cloud Foundry](CF.md)
* [CF MySQL](CF_MYSQL.md) brokered service
* [Kubo/Kubernetes](KUBO.md)
