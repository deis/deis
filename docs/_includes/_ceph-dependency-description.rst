The Deis :ref:`control-plane` makes use of `Ceph`_ to provide persistent storage for
the :ref:`registry`, :ref:`database`, and :ref:`logger` components. The additional
operational complexity of Ceph is tolerated because of the need for persistent
storage for platform high availability.

Alternatively, persistent storage can be achieved by running an external S3-compatible
blob store, PostgreSQL database, and log service. For users on AWS, the convenience
of Amazon S3 and Amazon RDS make the prospect of running a Ceph-less Deis cluster
quite reasonable.

Running a Deis cluster without Ceph provides several advantages:

* Removal of state from the control plane (etcd is still used for configuration)
* Reduced resource usage (Ceph can use up to 2GB of RAM per host)
* Reduced complexity and operational burden of managing Deis

.. _`Ceph`: http://ceph.com/
