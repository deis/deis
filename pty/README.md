pty
=========

container to run on-off commands. The container wraps the executed command and
allows the access stdin and stdout of the command via tcp.

## Running this Container

Set $COMMAND to the command you wish to execute:

    ><> docker run -d -v /var/run/docker.sock:/tmp/docker.sock -e COMMAND=ls deis/pty

## Building from Source

To build the image, run `make build`.
