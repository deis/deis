A discovery URL links `etcd`_ instances together by storing their peer
addresses and metadata under a unique identifier. Run this command from the root
of the repository to generate a ``contrib/coreos/user-data`` file with a new
discovery URL:

.. code-block:: console

    $ make discovery-url

Required scripts are supplied in this ``user-data`` file, so do not provision a
Deis cluster without running ``make discovery-url``.

.. _`etcd`: https://github.com/coreos/etcd
