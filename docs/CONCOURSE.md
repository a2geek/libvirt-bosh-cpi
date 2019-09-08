# Concourse

Set `CONCOURSE_DIR` to a local copy of the [Concourse BOSH deployment](https://github.com/concourse/concourse-bosh-deployment) directory.

Deploy!
```
$ bosh -n -d concourse deploy $CONCOURSE_DIR/cluster/concourse.yml \
    -o $CONCOURSE_DIR/cluster/operations/basic-auth.yml \
    -o $CONCOURSE_DIR/cluster/operations/static-web.yml \
    -o $CONCOURSE_DIR/cluster/operations/privileged-http.yml \
    -l $CONCOURSE_DIR/versions.yml \
    --vars-store=concourse-creds.yml \
    -l manifests/concourse-vars.yml
```

Concourse will be available at http://192.168.123.250 (assuming all the network stuff fits with your setup).
