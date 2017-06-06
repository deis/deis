# Provision a Deis Cluster on OpenStack

Please refer to the instructions at http://docs.deis.io/en/latest/installing_deis/openstack/.

Provision a cluster with a command such as:

```
$ make discovery-url
$ terraform apply \
      -var "username=$OS_USERNAME" \
      -var "password=$OS_PASSWORD" \
      -var "tenant=$OS_TENANT_NAME" \
      -var "auth_url=$OS_AUTH_URL" \
      -var "image=coreos-766.5.0" \
      contrib/openstack
```                  
