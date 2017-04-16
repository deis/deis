:title: Deploying Docker Images on Deis
:description: How to deploy applications on Deis using Docker Images

.. _using-docker-images:

Using Docker Images
===================
Deis supports deploying applications via an existing `Docker Image`_.
This is useful for integrating Deis into Docker-based CI/CD pipelines.

Prepare an Application
----------------------
Start by cloning an example application:

.. code-block:: console

    $ git clone https://github.com/deis/example-go.git
    $ cd example-go
    $ git checkout docker

Next use your local ``docker`` client to build the image and push
it to `DockerHub`_.

.. code-block:: console

    $ docker build -t <username>/example-go .
    $ docker push <username>/example-go

Docker Image Requirements
^^^^^^^^^^^^^^^^^^^^^^^^^
In order to deploy Docker images, they must conform to the following requirements:

* The Docker image must EXPOSE only one port
* The port must be listening for a HTTP connection
* A default CMD must be specified for running the container
* Bash should be available in the Docker image

.. note::

    Docker images which expose more than one port will hit `issue 1156`_.

Create an Application
---------------------
Use ``deis create`` to create an application on the :ref:`controller`.

.. code-block:: console

    $ mkdir -p /tmp/example-go && cd /tmp/example-go
    $ deis create
    Creating application... done, created example-go

.. note::

    The ``deis`` client uses the name of the current directory as the
    default app name.

Deploy the Application
----------------------
Use ``deis pull`` to deploy your application from `DockerHub`_ or
a private registry.

.. code-block:: console

    $ deis pull gabrtv/example-go:latest
    Creating build...  done, v2

    $ curl -s http://example-go.local3.deisapp.com
    Powered by Deis

Because you are deploying a Docker image, the ``cmd`` process type is automatically scaled to 1 on first deploy.

Use ``deis scale cmd=3`` to increase ``cmd`` processes to 3, for example. Scaling a
process type directly changes the number of :ref:`Containers <container>`
running that process.


.. attention::

    Support for Docker registry authentication is coming soon


.. _`Docker Image`: https://docs.docker.com/introduction/understanding-docker/
.. _`DockerHub`: https://registry.hub.docker.com/
.. _`CMD instruction`: https://docs.docker.com/reference/builder/#cmd
.. _`issue 1156`: https://github.com/deis/deis/issues/1156


Building Local Images with Deis/CoreOS
=======================================

When you start out with writing Dockerfiles, you need a place to test them - the
local Deis cluster can be used for that. However, the default setup of Deis aims
for security-by-default, so you'll have to prepare your cluster for development.

Therefore you'll have to uncomment a bit of the ``user-data`` file. Make sure
that ``contrib/coreos/user-data`` (assuming you've done ``make discovery-url``)
is uncommented, to look like this. (It's just after the fleet unit)

.. code-block:: yaml

  - name: docker.socket¬
    command: start¬
    drop-ins:¬
    - name: 30-ListenStream.conf¬
      content: |¬
        [Socket]¬
        ListenStream=2375¬

The cloud-config is updated every time you restart the node. The etcd cluster
supports being restarted one-at-a-time (if you restart more in a single go,
you're in for some pain, with the current CoreOS release).

.. code-block:: console

  vagrant reload deis-01

This will have updated the CloudInit (user-data) file inside deis-01. Now SSH
into it and run the following command (assuming you have a NFS mount point
./share):

.. code-block:: console

  sudo coreos-cloudinit --from-file share/contrib/coreos/user-data

This will list everything it's doing, and will look something like:

  Checking availability of "local-file"
  Fetching user-data from datasource of type "local-file"
  Fetching meta-data from datasource of type "local-file"
  2015/03/09 16:57:56 Parsing user-data as cloud-config
  Processing cloud-config from user-data
  2015/03/09 16:57:56 Writing file to "/etc/deis-release"
  2015/03/09 16:57:56 Wrote file to "/etc/deis-release"
  2015/03/09 16:57:56 Wrote file /etc/deis-release to filesystem
  // snip
  2015/03/09 16:57:56 Writing drop-in unit "30-ListenStream.conf" to filesystem
  2015/03/09 16:57:56 Writing file to "/etc/systemd/system/docker.socket.d/30-ListenStream.conf"
  2015/03/09 16:57:56 Wrote file to "/etc/systemd/system/docker.socket.d/30-ListenStream.conf"
  // snip

After updating the systemd configuration with a DropIn, you need to stop docker
and its socket, and then restart it again.

  systemctl stop docker
  systemctl stop docker.socket
  systemctl start docker.socket
  systemctl start docker

If you, outside the VM, ensure you have:

  export DOCKER_HOST=tcp://172.17.8.100:2375

You should now be able to talk with the docker daemon. Try it out!

If you get that your docker client doesn't match, do something like this on OS
X:

  cd $(brew --prefix)/Library/Formula
  git log -- docker.rb # find the right version from the build bot
  git show <sha1>:./docker.rb >docker.rb
  git status # verify that docker.rb is changed
  brew reinstall docker

Try using ``docker ps`` again!

If you find that the restart of docker shut some deis services down (docker ps
on the machine of your choice doesn't display as many running units), then you
can use the comand

  fleetctl list-units

It will show what machines have failed units -- you can manually start the unit
again -- ssh into the machine running docker and run

  sudo systemctl start deis-publisher.service && \
  journalctl -fu deis-publisher.service

...if it was deis-publisher which hadn't started.

Appendix - mounting NFS in docker
---------------------------------

Some docker features for building dockerfiles won't work unless you have the
files inside the virtual machine -- or perhaps you want to expose ``user-data``
inside the VM. You can add a line like this to ``Vagrantfile``:

  config.vm.synced_folder "#{ENV['HOME']}/dev/deis", "/home/core/share",
                          id: "ops", :nfs => true,
                          :mount_options => ['nolock,vers=3,udp']
