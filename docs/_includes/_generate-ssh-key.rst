The ``deisctl`` utility communicates with remote machines over an SSH tunnel.
If you don't already have an SSH key, the following command will generate
a new keypair named "deis":

.. code-block:: console

    $ ssh-keygen -q -t rsa -f ~/.ssh/deis -N '' -C deis
