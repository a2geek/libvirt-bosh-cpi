# Kubo (Kubernetes)

> NOTE: The CPI appears to handle to deployment without issues. However, the Kubernetes installation doesn't appear to handle storage volumes. It is uncertain if this is due to a lack of understanding on how to set it up, lack of Libvirt support, or possibly something related to the BOSH CPI (likely doubtful).

Generally, the Kubo Release instructions were used:

1. Install `kubectl` and `helm`
   * [`kubectl`](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
   * [`helm`](https://github.com/helm/helm/blob/master/docs/install.md)

2. Kubo is setup to deploy explicitly to the Ubuntu Xenial 315.41 stemcell. Upload that specific version:
   ```
   $ bosh upload-stemcell --sha1 87880fc81654ff2b377592fb53c377ead418a908 \
          https://s3.amazonaws.com/bosh-core-stemcells/315.41/bosh-stemcell-315.41-openstack-kvm-ubuntu-xenial-go_agent.tgz
   ```

3. Deploy Kubo:
   ```
   $ bosh deploy -n -d cfcr ${KUBO_DEPLOYMENT}/manifests/cfcr.yml \
       -o ${KUBO_DEPLOYMENT}/manifests/ops-files/misc/single-master.yml \
       -o ${KUBO_DEPLOYMENT}/manifests/ops-files/add-hostname-to-master-certificate.yml \
       -v api-hostname=kubo.greene.lan
   ```

4. Run `apply-specs` errand:
   ```
   $ bosh -d cfcr run-errand apply-specs
   ```

5. Run `smoke-tests` errand:
   ```
   $ bosh -d cfcr run-errand smoke-tests
   ```

6. Pull the CredHub secrets and sign into CredHub:
   ```
   $ source scripts/credhub-env.sh 
   $ credhub login
   Setting the target url: https://192.168.123.7:8844
   Login Successful
   ```

7. Hack `/etc/hosts` since the IP is not currently static:
   ```
   $ bosh -d cfcr vms
   Using environment '192.168.123.7' as client 'admin'

   Task 87. Done

   Deployment 'cfcr'

   Instance                                     Process State  AZ  IPs             VM CID                                   VM Type        Active  
   master/1ef6aeca-40bf-4bf4-bc9f-b6180807993d  running        z1  192.168.123.11  vm-07b6497c-7955-425f-50a3-1e962861b335  small          true  
   worker/c2d85266-92b3-4f6c-aeaa-ca410ca8b418  running        z2  192.168.123.13  vm-415bc2a0-4f98-4e16-7129-afedc621fef0  small-highmem  true  
   worker/f58c9dab-d6ce-429a-ad55-05ebf30f05cd  running        z3  192.168.123.14  vm-7e9d72e7-8047-4801-43a3-9c09a078c61c  small-highmem  true  
   worker/fb6d4270-fb65-4c77-8c74-bd94bee312cb  running        z1  192.168.123.12  vm-05456a1a-a4b1-4487-7942-e5d84cd48669  small-highmem  true  

   4 vms

   Succeeded
   ```
   Use the IP address for the `master` node (`192.168.123.11` in this sample).  Add that IP into your `/etc/hosts`:
   ```
   $ cat /etc/hosts | grep kubo
   192.168.123.11 kubo.greene.lan
   ```

8. Setup `kubectl`:
   ```
   $ ./bin/set_kubeconfig libvirt/cfcr https://kubo.greene.lan:8443
   Cluster "cfcr/libvirt/cfcr" set.
   User "cfcr/libvirt/cfcr/cfcr-admin" set.
   Context "cfcr/libvirt/cfcr" modified.
   Switched to context "cfcr/libvirt/cfcr".
   Created new kubectl context cfcr/libvirt/cfcr
   Try running the following command:
     kubectl get pods --namespace=kube-system
   ```

9. Setup `helm`:
   ```
   $ helm init
   $HELM_HOME has been configured at /home/rob/.helm.

   Tiller (the Helm server-side component) has been installed into your Kubernetes Cluster.

   Please note: by default, Tiller is deployed with an insecure 'allow unauthenticated users' policy.
   To prevent this, run `helm init` with the --tiller-tls-verify flag.
   For more information on securing your installation see: https://docs.helm.sh/using_helm/#securing-your-helm-installation
   ```

10. At this point, various errors were received (sample here):
    ```
    Error: configmaps is forbidden: User "system:serviceaccount:kube-system:default" cannot list resource "configmaps" in API group "" in the namespace "kube-system"
    ```
    So these actions were taken to grant access:
    ```
    $ kubectl create serviceaccount tiller --namespace kube-system
    serviceaccount/tiller created
    $ kubectl create -f manifests/tiller-clusterrolebinding.yml
    clusterrolebinding.rbac.authorization.k8s.io/tiller-clusterrolebinding created
    $ helm init --service-account tiller --upgrade
    $HELM_HOME has been configured at /home/rob/.helm.

    Tiller (the Helm server-side component) has been upgraded to the current version.
    ```

11. Try installing MySQL:
    ```
    $ helm install stable/mysql
    ```
    Observe that it does not get the persistent volume...
    ```
    $ kubectl describe pod kissing-elk-mysql-7d5cffc696-hbt5t
    Name:           kissing-elk-mysql-7d5cffc696-hbt5t
    <snip>
    Volumes:
      data:
        Type:       PersistentVolumeClaim (a reference to a PersistentVolumeClaim in the same namespace)
        ClaimName:  kissing-elk-mysql
        ReadOnly:   false
    <snip>
    Events:
      Type     Reason            Age                   From               Message
      ----     ------            ----                  ----               -------
      Warning  FailedScheduling  3m13s (x12 over 17m)  default-scheduler  pod has unbound immediate PersistentVolumeClaims (repeated 3 times)
      Warning  FailedScheduling  2m28s (x16 over 15m)  default-scheduler  pod has unbound immediate PersistentVolumeClaims (repeated 2 times)
    ```

:-(

# References

* [Kubo Release deployment instructions](https://github.com/cloudfoundry-incubator/kubo-release)
* [Kubo Deployment](https://github.com/cloudfoundry-incubator/kubo-deployment)
* [Install and Setup kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* [Installing Helm](https://github.com/helm/helm/blob/master/docs/install.md)
* [Kubernetes RBAC](https://docs.bitnami.com/kubernetes/how-to/configure-rbac-in-your-kubernetes-cluster/#use-case-2-enable-helm-in-your-cluster)
* [Kubernetes PersistentVolumes documentation](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)
* [GitLab Helm Deployment Guide](https://docs.gitlab.com/charts/installation/deployment.html)