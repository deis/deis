Private Docker registry
=======================

This directory provides a Vagrantfile and user-data file to provision and configure a CoreOS machine
which runs a private Docker registry. This is useful for testing Deis because it is significantly
faster than the public Docker registry.

To run the registry, in this directory simply:
```console
$ vagrant up
```

The registry will then be accessible at `172.21.12.100:5000` from any other local VM, including
Deis machines.
