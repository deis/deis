:title: Using a docker hub mirror
:description: Using a docker hub mirror for deis.

.. _using_a_docker_hub_mirror:

Using a docker hub mirror
=========================

.. important::

    Using a docker hub mirror is not recommended when your deis cluster is
    running in a public cloud or has a reliable and fast connection to the
    Docker Hub as it introduces another point of failure

When your connection to the Docker Hub is limited, it might make sense to
to setup a hub mirror in your own network. This ensures that docker layers are
only downloaded once when required by multiple instances and allows you to
pre-warm your mirror, when setting up your deis cluster is timing out, because
of your slow connection.

Installing a local hub mirror
-----------------------------
Follow the step 2 of the `guide`_ provided by Docker to setup your local hub mirror. All deis
instances need access to this registry and you can not use the internal deis-registry
for this purpose.

Using your mirror
-----------------
After your local mirror is started, you need to add

.. code-block:: console

    --registry-mirror=http://YOUR_IP:YOUR_PORT

to ``DOCKER_OPTS`` in your `contrib/coreos/user-data`_.

.. _`guide`: https://docs.docker.com/articles/registry_mirror/#step-2-run-the-local-registry-mirror
.. _`contrib/coreos/user-data`: https://github.com/deis/deis/blob/master/contrib/coreos/user-data.example
