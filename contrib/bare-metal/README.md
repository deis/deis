# Provision a Deis Cluster on bare-metal hardware

Deis clusters can be provisioned anywhere [CoreOS](https://coreos.com/) can, including on your own hardware. To get CoreOS running on raw hardware, you can boot with [PXE](https://coreos.com/docs/running-coreos/bare-metal/booting-with-pxe/) or [iPXE](https://coreos.com/docs/running-coreos/bare-metal/booting-with-ipxe/) - this will boot a CoreOS machine running entirely from RAM. Then, you can [install CoreOS to disk](https://coreos.com/docs/running-coreos/bare-metal/installing-to-disk/).

## Generate SSH key
To avoid problems deploying/launching apps later on it is necessary to install [CoreOS](https://coreos.com/) to disk with a SSH key without a passphrase. The following command will generate a new keypair named "deis".

```console
$ ssh-keygen -q -t rsa -f ~/.ssh/deis -N '' -C deis
```

## Customize user-data

### Discovery URL
Copy [user-data.example](../coreos/user-data.example) to [user-data](../coreos/user-data)
and add a new discovery URL.
You can get a new one by sending a request to http://discovery.etcd.io/new.

### SSH Key
Add the public key part for the SSH key generated in the first step to the [user-data](../coreos/user-data.example) file:

```yaml
ssh_authorized_keys:
  - ssh-rsa AAAAB3... deis
```

### Update $private_ipv4
[CoreOS](https://coreos.com/) on bare metal doesn't detect the `$private_ipv4` reliably. Replace all occurences in the [user-data](../coreos/user-data.example) with the (private) IP address of the node.

### Add environment
Since [CoreOS](https://coreos.com/) doesn't detect private and public IP adresses the `/etc/environmnet` file doesn't get written on boot. Add it to the `write_files` section of [user-data](../coreos/user-data.example)

```yaml
  - path: /etc/environment
    permissions: 0644
    content: |
      COREOS_PUBLIC_IPV4=<your public ip>
      COREOS_PRIVATE_IPV4=<your private ip>
```

## Install CoreOS to disk
Assuming you have booted your bare metal server into [CoreOS](https://coreos.com/) you can perform now perform the installation to disk.

### Provide the config file to the installer
Save the [user-data](../coreos/user-data.example) to your bare metal machine. The example assumes you transferred the config to `/tmp/config`

### Start the installation
```console
coreos-install -C alpha -c /tmp/config -d /dev/sda
```

This will install the current [CoreOS](https://coreos.com/) release to disk. If you want to install the recommended [CoreOS](https://coreos.com/) version check the [Deis changelog](../../CHANGELOG.md)
and specify that version by appending the `-V` parameter to the install command, e.g. `-V 402.2.0`.

After the installation has finished reboot your server. Once your machine is back up you should be able to log in as the `core` user using the `deis` ssh key.

## Initialize the cluster
Once your server(s) are all provisioned you can proceed to install Deis. Use the hostname of one of your machines in the next step.

```console
$ ssh-add ~/.ssh/deis
$ export FLEETCTL_TUNNEL=your.server.name.here
$ cd ../.. && make run
```

## Use Deis!
After that, register with Deis!
```console
$ deis register http://deis.example.org
username: deis
password:
password (confirm):
email: info@opdemand.com
```

## Considerations when deploying Deis:
* Use machines with ample disk space and RAM (we use [large instances](https://aws.amazon.com/ec2/instance-types/) on EC2, for comparison)
* Choose an appropriate [cluster size](https://github.com/coreos/etcd/blob/master/Documentation/optimal-cluster-size.md)

## Known problems

### Hostname is localhost
If your hostname after installation to disk is `localhost` set the hostname in [user-data](../coreos/user-data.example) before installation:

```yaml
hostname: your-hostname
```

The hostname must not be the fully qualified domain name!

### Slow name resolution

Certain DNS servers and firewalls have problems with glibc sending out requests for IPv4 and IPv6 addresses in parallel. The solution is to set the option `single-request` in `/etc/resolv.conf`. This can best be accomplished in the [user-data](../coreos/user-data.example) when installing [CoreOS](https://coreos.com/) to disk. Add the following block to the `write_files` section:

```yaml
  - path: /etc/resolv.conf
    permissions: 0644
    content: |
      nameserver 8.8.8.8
      nameserver 8.8.4.4
      domain your.domain.name
      options single-request
```
