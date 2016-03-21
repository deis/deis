The `source code`_ for Deis must be on your workstation to run the commands in
this documentation. Download an archive file from the `releases page`_, or use
``git`` to clone the repository:

.. code-block:: console

    $ git clone https://github.com/deis/deis.git
    $ cd deis
    $ git checkout v1.13.0

Check out the latest Deis release, rather than using the default (master).

If you contribute to Deis or build components locally, use ``go get`` instead to
clone the source code into your `$GOPATH`_:

.. code-block:: console

    $ go get -u -v github.com/deis/deis
    $ cd $GOPATH/src/github.com/deis/deis

Additionally, you'll need the ``deisctl`` CLI tool. If you don't already have it,
install instructions are :ref:`here <install_deisctl>`.

.. _`source code`: https://github.com/deis/deis
.. _`releases page`: https://github.com/deis/deis/releases
.. _`$GOPATH`: http://golang.org/doc/code.html#GOPATH
