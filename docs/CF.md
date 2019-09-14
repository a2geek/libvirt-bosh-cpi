# Cloud Foundry

Set `CF_DEPLOYMENT_DIR` to a local copy of the [Cloud Foundry deployment](https://github.com/cloudfoundry/cf-deployment/) directory.

Add the DNS runtime config:
```
$ bosh update-runtime-config $BOSH_DEPLOYMENT_DIR/runtime-configs/dns.yml --name dns
```

Deploy!
```
$ bosh -n -d cf deploy $CF_DEPLOYMENT_DIR/cf-deployment.yml \
    -o $CF_DEPLOYMENT_DIR/operations/scale-to-one-az.yml \
    -o $CF_DEPLOYMENT_DIR/operations/set-router-static-ips.yml \
    -o $CF_DEPLOYMENT_DIR/operations/use-compiled-releases.yml \
    -o $CF_DEPLOYMENT_DIR/operations/use-latest-stemcell.yml \
    -o $CF_DEPLOYMENT_DIR/operations/override-app-domains.yml \
    -l manifests/cloudfoundry-vars.yml \
    --vars-store=cf-creds.yml
```

Note that there is a requirement for DNS resolution to `*.sys.mycf.lan` as currently configured. `/etc/hosts` can be used as a hack for validation.

```
$ cat /etc/hosts | grep mycf
192.168.123.252 api.sys.mycf.lan login.sys.mycf.lan sample1.sys.mycf.lan
```

> If you use DD-WRT, you can also add wildcard entries of `address=/.mycf.lan/192.168.123.252` and `address=/.sys.mycf.lan/192.168.123.252` and setup a routing entry to direct all `192.168.123.*` entries to the host machine's IP address.

To get the CF admin credentials, there are a few hoops.

Get the CredHub admin credentials and login to CredHub:
```
$ source scripts/credhub-env.sh 
$ credhub login
Setting the target url: https://192.168.123.7:8844
Login Successful
```

> Note that the next sequence assumes you have a fairly recent version of `cf` installed (mine is currently at 
6.46.1+4934877ec.2019-08-23).

To get the CF admin credentials and login to CF _as an admin_:
```
$ source scripts/cf-env.sh
$ cf api https://api.sys.mycf.lan --skip-ssl-validation
Setting api endpoint to https://api.sys.mycf.lan...
OK

api endpoint:   https://api.sys.mycf.lan
api version:    2.139.0
$ cf auth
API endpoint: https://api.sys.mycf.lan
Authenticating...
OK

Use 'cf target' to view or set your target org and space.
```

Finally, create a place to deploy applications:
```
$ cf create-org rob
Creating org rob as admin...
OK

Assigning role OrgManager to user admin in org rob...
OK

TIP: Use 'cf target -o "rob"' to target new org
$ cf target -o rob
api endpoint:   https://api.sys.mycf.lan
api version:    2.139.0
user:           admin
org:            rob
No space targeted, use 'cf target -s SPACE'
$ cf create-space dev
Creating space dev in org rob as admin...
OK

Assigning role SpaceManager to user admin in org rob / space dev as admin...
OK

Assigning role SpaceDeveloper to user admin in org rob / space dev as admin...
OK

TIP: Use 'cf target -o "rob" -s "dev"' to target new space
$ cf target -s dev
api endpoint:   https://api.sys.mycf.lan
api version:    2.139.0
user:           admin
org:            rob
space:          dev
```

Note that `sample1.sys.mycf.lan` is just for a quick test deploy like this:
```
$ mkdir staticfile-sample
$ cd staticfile-sample
staticfile-sample$ cat > index.html
Hello World!
^D
staticfile-sample$ cat > manifest.yml
applications:
- name: sample1
  buildpacks:
  - staticfile_buildpack
  memory: 32M
  routes:
  - route: sample1.mycf.lan
^D
staticfile-samplestatic-site$ cf push -f manifest.yml -p .
<snip>
staticfile-sample$ cf apps
Getting apps in org rob / space dev as admin...
OK

name      requested state   instances   memory   disk   urls
sample1   started           1/1         32M      1G     sample1.mycf.lan
$ curl http://sample1.mycf.lan
Hello World!
```

# References

* [Cloud Foundry Deployment Guide](https://github.com/cloudfoundry/cf-deployment/blob/master/texts/deployment-guide.md)
* [Cloud Foundry Cloud Configs](https://github.com/cloudfoundry/cf-deployment/blob/master/texts/on-cloud-configs.md)
