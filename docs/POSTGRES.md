# Postgres

Set `POSTGRES_DIR` to a local copy of the [Postgres release](https://github.com/cloudfoundry/postgres-release) directory.

Upload a Postgres release...
```
$ bosh upload-release https://bosh.io/d/github.com/cloudfoundry/postgres-release
```

Deploy!
```
$ bosh -n -d postgres deploy $POSTGRES_DIR/templates/postgres.yml \
    -o $POSTGRES_DIR/templates/operations/add_static_ips.yml \
    -o $POSTGRES_DIR/templates/operations/set_properties.yml \
    -o $POSTGRES_DIR/templates/operations/use_bbr.yml \
    --vars-store=postgres-creds.yml \
    -l manifests/postgres-vars.yml
```
