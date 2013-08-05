:title: Key
:description: What is a Deis Key?
:keywords: key, ssh

.. _key:

Key
===
Deis keys are SSH Keys used during the git push process.  Each user
can use the client to manage a list of keys on the :ref:`Controller`.

A user's keys are automatically added to every launched :ref:`Node`,
allowing easy SSH access via:

.. code-block:: console

	$ ssh ubuntu@<node-fqdn>
