:title: Recovering Ceph quorum
:description: Additional information for recovering clusters once Ceph has lost quorum.

.. _recovering-ceph-quorum:

Recovering Ceph quorum
======================

Ceph relies on `Paxos`_ to maintain a quorum among monitor services so that they agree on cluster state.
In some cases Ceph can lose quorum, such as when hosts are added and removed from the cluster in
quick succession, without removing the old hosts from Ceph (see :ref:`add_remove_host`).

A telltale sign of quorum loss is when querying cluster health, ``ceph -s`` times out with monitor
faults on every host in the cluster.

.. important::

    Ceph refusing to do anything when it has lost quorum is a safety precaution to prevent you
    from losing data. Attempting to recover from this situation requires knowledge about the state
    of your cluster, and should only be attempted if data loss is not considered catastrophic (such as
    when a recent backup is available). When in doubt, consult the Ceph and Deis communities for
    assistance. Deis recommends regular backups to minimize impact should an issue like this occur.
    For more information, see :ref:`backing_up_data`.

The instructions below are intentionally vague, as each recovery scenario will be unique. They are
intended only to point users in the right direction for recovery.

To recover from Ceph quorum loss:

#. Suspect quorum loss because ``ceph -s`` shows nothing but timeouts and/or monitor faults
#. :ref:`using-store-admin`, use the Ceph `admin socket`_ to query the `mon status`_, identifying that there are enough stale entries to prevent Ceph from gaining quorum
#. Stop the platform with ``deisctl stop platform`` so components stop trying to write data to store (note that instead, manually stopping all components except router will allow application containers to remain up, unaffected)
#. Clean up stale entries in ``/deis/store/hosts`` so that dead monitors are not written out to clients
#. Update ``/deis/store/monSetupLock`` to point to the healthy monitor -- note that this isn't strictly necessary, as this value is only used if wiping clean and starting a fresh cluster from scratch with no data, but it's good cleanup
#. Start the healthy monitor and use the admin socket to get the current state of the cluster.
#. Given the cluster state as the monitor sees it, use `monmaptool`_ to manually remove stale monitor entries from the monmap (i.e. ``monmaptool --rm mon.<hostname> --clobber /etc/ceph/monmap``)
#. Stop the healthy monitor and use ``deis-store-admin`` to inject the prepared monmap into the monitor with ``ceph-mon -i <hostname> --inject-monmap /etc/ceph/monmap``
#. Start the monitor and ensure it achieves quorum by itself (use ``ceph -s`` and/or query mon_status on the admin socket)
#. Start the other monitors and ensure they connect
#. Start the OSDs with ``deisctl start store-daemon``
#. Observe the OSD map with ``ceph osd dump`` -- for each OSD that is no longer with us, follow :ref:`removing_an_osd` -- take care to ensure that the data is relocated (watch the health with ``ceph -w``) before marking another OSD as ``out``
#. Once the OSD map reflects the now-healthy OSDs, start the remaining store services in order: ``deisctl start store-metadata`` and ``deisctl start store-gateway``
#. Confirm that the cluster is healthy with the metadata servers added, and then start ``store-volume`` with ``deisctl start store-volume``.
#. Start the remaining services with ``deisctl start platform``

.. _`admin socket`: http://ceph.com/docs/master/rados/troubleshooting/troubleshooting-mon/#using-the-monitor-s-admin-socket
.. _`mon status`: http://ceph.com/docs/master/rados/troubleshooting/troubleshooting-mon/#understanding-mon-status
.. _`monmaptool`: http://ceph.com/docs/master/man/8/monmaptool/
.. _`Paxos`: http://en.wikipedia.org/wiki/Paxos_%28computer_science%29
