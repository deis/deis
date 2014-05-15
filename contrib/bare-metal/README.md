# Provision a Deis Cluster on bare-metal hardware

Deis clusters can be provisioned anywhere [CoreOS](https://coreos.com/) can, including on your own hardware. To get CoreOS running on raw hardware, you can boot with [PXE](https://coreos.com/docs/running-coreos/bare-metal/booting-with-pxe/) or [iPXE](https://coreos.com/docs/running-coreos/bare-metal/booting-with-ipxe/) - this will boot a CoreOS machine running entirely from RAM. Then, you can [install CoreOS to disk](https://coreos.com/docs/running-coreos/bare-metal/installing-to-disk/).

Considerations when deploying Deis:
* Use machines with ample disk space and RAM (we use [large instances](https://aws.amazon.com/ec2/instance-types/) on EC2, for comparison)
* Choose an appropriate [cluster size](https://github.com/coreos/etcd/blob/master/Documentation/optimal-cluster-size.md)
* Supply our [cloud config file](../coreos/user-data), making sure to use a [new discovery URL](https://discovery.etcd.io/new)
* Use the `alpha` channel of CoreOS

We hope to improve our documentation around bare metal provisioning. If you're deployed Deis on bare metal and think you can help improve this documentation, please submit a pull request. Thanks!
