:title: Isolating the Planes
:description: Configuring the cluster to isolate the control plane, data plane, and router mesh.

.. _isolating-planes:

Isolating the Planes
====================

.. include:: ../_includes/_isolating-planes-description.rst

Understanding Fleet metadata
----------------------------

The key to isolating the Control Plane, Data Plane, and Router Mesh is Fleet
metadata.  Although Deis supports alternate schedulers, Deis components
themselves are all scheduled via Fleet.

Deis configures the Fleet daemon executing on each node at the time of
provisioning via cloud-config.  Within that configuration, it is possible to tag
nodes with metadata in the form of key/value pairs to arbitrarily describe
attributes of the node.  For instance, an operator may tag a node with
``ssd=true`` to indicate that a node's volumes use solid state disk.

.. code-block:: yaml

    #cloud-config
    ---
    coreos:
      fleet:
        metadata: ssd=true
    # ...

When scheduling a unit of work via Fleet, it is also possible to annotate that
unit with metadata that is required to be present on any node in order to be
considered eligible to host that work.  In keeping with our previous example,
to restrict a unit of work to only those nodes equipped with SSD, the unit may
be annotated thusly:

.. code-block:: yaml

    # ...
    [X-Fleet]
    MachineMetadata="ssd=true"

Deis takes advantage of this very mechanism to establish which nodes are
eligible to host each of the Control Plane, Data Plane, and Router Mesh.

`More details on Fleet metadata`_

cloud-config
------------

To configure a Fleet node as eligible to host Control Plane components, the
following cloud-config may be used:

.. code-block:: yaml

    #cloud-config
    ---
    coreos:
      fleet:
        metadata: controlPlane=true

Similarly, ``dataPlane=true`` and ``routerMesh=true`` may be used to establish
eligibility to host components of the Data Plane (including applications) and
Router Mesh, respectively.

It is also possible to configure nodes as eligible to host two or even all
three of the Control Plane, Data Plane, and Router Mesh.  In fact, this is
the default behavior described by Deis' included cloud-config.

.. code-block:: yaml

    #cloud-config
    ---
    coreos:
      fleet:
        metadata: controlPlane=true,dataPlane=true,routerMesh=true

It should be obvious that isolating the planes as described here requires
subsets of a cluster's nodes to be configured differently from one another (with
different metadata). Deis provisioning scripts do not currently account for
this, so managing separate cloud-config for each subset of nodes in the cluster
is left as an exercise for the advanced operator.

Decorating units
----------------

To complement the cloud-config described above, Deis 1.10.0 and later are capable
of seamlessly "decorating" the Fleet units for each Deis platform component with
the metadata that describes where each unit may be hosted.

.. note::

    For the purposes of backwards compatibility with Deis clusters provisioned
    using versions of Deis older than 1.10.0, decorating the platform's units
    with metadata is an opt-in.  Nodes in older clusters are guaranteed to be
    lacking the metadata that indicates what components they are eligible to
    host.  As such, decorated units would be ineligible to run anywhere within
    such a cluster.

    To opt in, use the following:

    .. code-block:: console

        $ deisctl config platform set enablePlacementOptions=true


.. _`More details on Fleet metadata`: https://coreos.com/fleet/docs/latest/unit-files-and-scheduling.html#fleet-specific-options
