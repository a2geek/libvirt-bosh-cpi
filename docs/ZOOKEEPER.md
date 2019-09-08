# Zookeeper

Set `ZOOKEEPER_DIR` to a local copy of the [Zookeeper release](https://github.com/cppforlife/zookeeper-release) directory.

Deploy!
```
$ bosh -n -d zookeeper deploy $ZOOKEEPER_DIR/manifests/zookeeper.yml \
    --vars-store=zookeeper-creds.yml
```
