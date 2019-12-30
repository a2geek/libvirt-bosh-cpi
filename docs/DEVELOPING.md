# Developing

At the command-line, from the source directory a compile can be done:

```
$ cd src
$ go build -mod=vendor -o a.out main/main.go
$ rm a.out
```

Also, tests can be run:

```
$ cd src
$ go test -v ./...
```

# Agent configuration 

The agent configuration structures are here:

* [BOSH Agent MetadataContentsType](https://godoc.org/github.com/cloudfoundry/bosh-agent/infrastructure#MetadataContentsType)
* [BOSH Agent UserDataContentsType](https://godoc.org/github.com/cloudfoundry/bosh-agent/infrastructure#UserDataContentsType)
