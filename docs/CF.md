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
    -l manifests/cloudfoundry-vars.yml \
    --vars-store=cf-creds.yml
```

Note that there is a requirement for DNS resolution to `*.sys.mypcf.lan` as currently configured. `/etc/hosts` can be used as a hack for validation.

```
$ cat /etc/hosts | grep mypcf
192.168.123.252 api.sys.mypcf.lan login.sys.mypcf.lan sample1.sys.mypcf.lan
```

To get the admin credentials setup:
```
$ export CREDHUB_CLIENT=credhub-admin
$ export CREDHUB_SECRET=$(bosh interpolate ./cf-creds.yml --path=/credhub_admin_client_secret)
$ export CREDHUB_CA_CERT="$(bosh interpolate ./cf-creds.yml --path=/credhub_tls/ca )"$'\n'"$( bosh interpolate ./cf-creds.yml --path=/uaa_ssl/ca)"
```

To login with those credentials:
```
$ cf api --skip-ssl-validation https://api.sys.mypcf.lan
Setting api endpoint to https://api.sys.mypcf.lan...
OK

api endpoint:   https://api.sys.mypcf.lan
api version:    2.139.0
$ cf login 
API endpoint: https://api.sys.mypcf.lan

Email> admin

Password> (paste in cf_admin_password from cf-creds.yml file)
Authenticating...
OK

Targeted org system

API endpoint:   https://api.sys.mypcf.lan (API version: 2.139.0)
User:           admin
Org:            system
Space:          No space targeted, use 'cf target -s SPACE'
```

Finally, create a place to deploy applications:
```
$ cf create-org robstuff
Creating org robstuff as admin...
OK

Assigning role OrgManager to user admin in org robstuff ...
OK

TIP: Use 'cf target -o "robstuff"' to target new org

$ cf target -o robstuff
api endpoint:   https://api.sys.mypcf.lan
api version:    2.139.0
user:           admin
org:            robstuff
No space targeted, use 'cf target -s SPACE'

$ cf create-space np
Creating space np in org robstuff as admin...
OK
Assigning role RoleSpaceManager to user admin in org robstuff / space np as admin...
OK
Assigning role RoleSpaceDeveloper to user admin in org robstuff / space np as admin...
OK

TIP: Use 'cf target -o "robstuff" -s "np"' to target new space

$ cf target -s np
api endpoint:   https://api.sys.mypcf.lan
api version:    2.139.0
user:           admin
org:            robstuff
space:          np
```

Note that `sample1.sys.mypcf.lan` is just for a quick test deploy like this:
```
$ mkdir staticfile-sample
$ cd staticfile-sample
staticfile-sample$ touch Staticfile
staticfile-sample$ cat > index.html
Hello World!
^D
staticfile-sample$ cf push -b staticfile_buildpack -m 32M -p . sample1
staticfile-sample$ cf apps
Getting apps in org robstuff / space np as admin...
OK

name      requested state   instances   memory   disk   urls
sample1   started           1/1         32M      1G     sample1.sys.mypcf.lan
```

# References

* [Cloud Foundry Deployment Guide](https://github.com/cloudfoundry/cf-deployment/blob/master/texts/deployment-guide.md)
* [Cloud Foundry Cloud Configs](https://github.com/cloudfoundry/cf-deployment/blob/master/texts/on-cloud-configs.md)
