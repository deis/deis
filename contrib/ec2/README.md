# Provision a Deis Cluster on Amazon EC2

## Install the [AWS Command Line Interface][aws-cli]:
```console
$ pip install awscli
Downloading/unpacking awscli
  Downloading awscli-1.3.6.tar.gz (173kB): 173kB downloaded
  ...
```

## Configure aws-cli
Run `aws configure` to set your AWS credentials:
```console
$ aws configure
AWS Access Key ID [None]: ***************
AWS Secret Access Key [None]: ************************
Default region name [None]: us-west-1
Default output format [None]:
```

## Upload keys
Upload a new keypair to AWS, ensuring that the name of the keypair is set to "deis".

## Customize cloudformation.json
Edit [cloudformation.json][cf-params], ensuring to add a new discovery URL.
You can get a new one by sending a new request to http://discovery.etcd.io/new.
```console
    {
        "ParameterKey":     "DiscoveryURL",
        "ParameterValue":   "https://discovery.etcd.io/40826e8da55f4d9026935ab67b243c6a"
    }
```
NOTE: If you're interested in running your own discovery endpoint or want to know more
about the discovery URL, see http://discovery.etcd.io for more information. You can also
read more on how you can customize this cluster by looking at the
[CoreOS EC2 template][template] and applying it to [cloudformation.json][cf-params].

## Run the provision script
Run the [cloudformation provision script][pro-script] to spawn a new CoreOS cluster:
```console
$ ./provision-ec2-cluster.sh
{
    "StackId": "arn:aws:cloudformation:us-west-1:413516094235:stack/deis/9699ec20-c257-11e3-99eb-50fa01cd4496"
}
Your Deis cluster has successfully deployed.
Please wait for it to come up, then run ./initialize-ec2-cluster.sh
```

## Initialize the cluster
Once the cluster is up, get the hostname of any of the machines from EC2, set
FLEETCTL_TUNNEL, then run [the init script][init-script] to bootstrap the cluster
remotely:
```console
$ ssh-add ~/.ssh/id_rsa
$ export FLEETCTL_TUNNEL=ec2-12-345-678-90.us-west-1.compute.amazonaws.com
$ ./initialize-ec2-cluster.sh
The authenticity of host '54.215.248.50:22' can't be established.
RSA key fingerprint is 86:10:74:b9:6a:ee:3b:21:d0:0f:b4:63:cc:10:64:c9.
Are you sure you want to continue connecting (yes/no)? yes
Warning: Permanently added '54.215.248.50:22' (RSA) to the list of known hosts.
Job deis-registry.service started on aec641dc.../172.31.21.4
Job deis-logger.service started on 494dcb6a.../172.31.5.226
Job deis-database.service started on aec641dc.../172.31.21.4
Job deis-cache.service started on aec641dc.../172.31.21.4
Job deis-controller.service started on aec641dc.../172.31.21.4
Job deis-builder.service started on 494dcb6a.../172.31.5.226
Job deis-router.service started on aec641dc.../172.31.21.4
done!
```

## Run Deis!
After that, wait for the components to come up, check which host the controller is
running on and register with Deis!
```
$ fleetctl list-units
UNIT                    LOAD    ACTIVE  SUB     DESC            MACHINE
deis-builder.service    loaded  active  running deis-builder    d9f1f3ea.../172.31.5.62
deis-cache.service      loaded  active  running deis-cache      d9f1f3ea.../172.31.5.62
deis-controller.service loaded  active  running deis-controller d9f1f3ea.../172.31.5.62
deis-database.service   loaded  active  running deis-database   13c5541b.../172.31.5.61
deis-logger.service     loaded  active  running deis-logger     d9f1f3ea.../172.31.5.62
deis-registry.service   loaded  active  running deis-registry   4c263e91.../172.31.24.155
deis-router.service     loaded  active  running deis-router     13c5541b.../172.31.5.61
$ deis register ec2-12-345-678-90.us-west-1.compute.amazonaws.com:8000
username: deis
password:
password (confirm):
email: info@opdemand.com
```

Follow the instructions in the [README](../../README.md) to deploy your first application, or check out the [docs](http://docs.deis.io/en/latest/) for more information. Have fun!

[aws-cli]: https://github.com/aws/aws-cli
[template]: https://s3.amazonaws.com/coreos.com/dist/aws/coreos-alpha.template
[cf-params]: cloudformation.json
[pro-script]: provision-ec2-cluster.sh
[init-script]: initialize-ec2-cluster.sh
