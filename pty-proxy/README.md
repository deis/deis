pty-proxy
=========

## Running this Container

Set $ETCD_HOST to be the IP address/hostname of the etcd endpoint you wish to target, and
$HOST to be the IP address of the host running this container:

    ><> docker run -d -e ETCD_HOST=192.168.0.1 -e HOST=192.168.0.1 deis/pty-proxy

## Building from Source

To build the image, run `make image`.

The build/runtime environment is split into two parts:

### The build environment

Based on deis/go, this image installs Go and compiles pty-proxy into a binary.
