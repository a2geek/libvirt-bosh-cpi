# MinIO


```
$ bosh upload-release https://bosh.io/d/github.com/minio/minio-boshrelease
```

```
$ bosh deploy -n -d minio manifests/minio-manifest.yml \
       -v minio_deployment_name=minio \
       --vars-store=minio-creds.yml
```

```
$ mc alias set tests3 \
    $(bosh int manifests/minio-vars.yml --path=/minio_dns) \
    $(bosh int minio-creds.yml --path=/minio_accesskey) \
    $(bosh int minio-creds.yml --path=/minio_secretkey)
```