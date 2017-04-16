pty
=========

## Running this Container

Set $ETCD_HOST to be the IP address/hostname of the etcd endpoint you wish to target, and
$HOST to be the IP address of the host running this container:

    ><> docker run -d -v /var/run/docker.sock:/tmp/docker.sock -e DOCKER_HOST=/tmp/docker.sock -e COMMAND=ls deis/pty
