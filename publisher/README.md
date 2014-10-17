publisher
=========

Publisher listens directly to a docker socket bind-mounted into the container and listens to the
docker events API for running containers on the host. Deis applications are published to etcd for
service discovery.

## Running this Container

Set $ETCD_HOST to be the IP address/hostname of the etcd endpoint you wish to target, and
$HOST to be the IP address of the host running this container:

    $ docker run -d -v /var/run/docker.sock:/tmp/docker.sock -e ETCD_HOST=192.168.0.1 -e HOST=192.168.0.1 deis/publisher

## Building from Source

To build the image, run `make build`.

The build/runtime environment is split into two parts:

### The build environment

Based on deis/go, this image installs Go and compiles publisher into a binary.

### The runtime environment

Leveraging the build environment, this image pulls in the standalone binary compiled in
the build environment and injects it into a minimal standalone container, minimizing the
disk space footprint that this image takes up. In fact, this image is < 5MB:

    $ docker images | grep publisher
    deis/publisher                           master              7974d140b07d        11 minutes ago      4.678 MB
    deis/publisher-build                     master              75983660e714        11 minutes ago      1.091 GB
