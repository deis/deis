:title: Key
:description: Deis keys are SSH keys used during the git push process. Use the Deis client to manage a list of keys on the Deis controller.
:keywords: key, ssh, deis

.. _key:

Key
===
Deis keys are SSH Keys used during the git push process.  Each user
can use the client to manage a list of keys on the :ref:`Controller`.

A user's keys are automatically added to every launched :ref:`Node`,
allowing easy SSH access via:

.. code-block:: console

	$ ssh ubuntu@<node-fqdn>
