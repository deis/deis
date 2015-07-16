:title: Disk usage
:description: Understanding disk usage for CoreOS and Deis.

.. _disk_usage:

Disk Usage
==========

When planning a Deis deployment, it's helpful to understand how both Deis and CoreOS utilize
local storage on a machine.

The following filesystem paths are the most important to consider:

===================      =============================        ============================================================================================================================================
location                 purpose                              considerations
===================      =============================        ============================================================================================================================================
/var/lib/etcd            etcd snapshot data                   etcd writes a relatively small amount of snapshot data here, so access should be as fast as possible (cloud providers use fast, local disks)
/var/lib/docker          Docker image/volume storage          should be large - on cloud providers with external storage (AWS, GCE, Azure) this is a separate 100GB volume
/var/lib/deis/store      mounted CephFS for deis-store        none - this is a virtually-mounted filesystem, and the "real" Ceph data lives in a Docker volume (so it's stored in /var/lib/docker)
/                        everything else (logs, etc.)         should be adequately large enough to prevent out-of-space issues causing service failure (on AWS this is a 50GB volume)
===================      =============================        ============================================================================================================================================

Identifying low disk space
~~~~~~~~~~~~~~~~~~~~~~~~~~

Usually, errors in component logs like "No space left on device" will clearly indicate that a
low disk space condition is the culprit of operational issues. Upon investigation, ``df -h`` should
reveal a filesystem with low free disk space.

In some cases, however, the output from ``df -h`` doesn't show any volume having low free space.
This typically points to a btrfs issue - see `btrfs troubleshooting`_ for more information.

Recovering disk space
~~~~~~~~~~~~~~~~~~~~~

If a volume is nearly full, it may be necessary to prune old data from it to ensure the cluster
remains operational.

The root volume should rarely become full. If it does, explore pruning old log files (or look for
a forgotten backup or download in the ``core`` user's home directory).

The most alarming low-disk-space condition is when the Docker volume is nearly full. The ``builder``
component should remove unnecessary images after a build, and will also remove images after
an application has been deleted.

However, it some cases it may be necessary to manually prune old images using ``docker rmi``:

.. code-block:: console

    $ docker images -aq | xargs -l10 docker rmi

.. note::

    This command actually instructs Docker to remove **all** images, and relies on the daemon's
    refusal to remove images which are in-use (errors will be emitted for running images).


.. _`btrfs troubleshooting`: https://coreos.com/docs/cluster-management/debugging/btrfs-troubleshooting/
