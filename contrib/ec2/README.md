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
Generate and upload a new keypair to AWS, ensuring that the name of the keypair is set to "deis".
```console
$ ssh-keygen -q -t rsa -f ~/.ssh/deis -N '' -C deis
$ aws ec2 import-key-pair --key-name deis --public-key-material file://~/.ssh/deis.pub
```

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
$ cd contrib/ec2
$ ./provision-ec2-cluster.sh
{
    "StackId": "arn:aws:cloudformation:us-west-1:413516094235:stack/deis/9699ec20-c257-11e3-99eb-50fa01cd4496"
}
Your Deis cluster has successfully deployed.
Please wait for all instances to come up as "running" before continuing.
```

## Initialize the cluster
Once the cluster is up, get the hostname of any of the machines from EC2, set
FLEETCTL_TUNNEL, and issue a `make run` from the project root:
```console
$ ssh-add ~/.ssh/deis
$ export FLEETCTL_TUNNEL=ec2-12-345-678-90.us-west-1.compute.amazonaws.com
$ cd ../.. && make run
```
The script will deploy Deis and make sure the services start properly.

## Configure DNS
While you can reference the controller and hosted applications with public hostnames provided by EC2, it is recommended for ease-of-use that
you configure your own DNS records using a domain you own. See [Configuring DNS](http://docs.deis.io/en/latest/operations/configure-dns/) for details.

## Use Deis!
After that, register with Deis!
```
$ deis register deis.example.org:8000
username: deis
password:
password (confirm):
email: info@opdemand.com
```

[aws-cli]: https://github.com/aws/aws-cli
[template]: https://s3.amazonaws.com/coreos.com/dist/aws/coreos-alpha.template
[cf-params]: cloudformation.json
[pro-script]: provision-ec2-cluster.sh
