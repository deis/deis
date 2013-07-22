deis
====

**Deis** is an open source *platform-as-a-service* (PaaS) for public and
private clouds.

Take your agile development to the next level. Free your mind and focus
on your code. Deploy updates to local metal or worldwide clouds with
``git push``. Scale servers, processes, and proxies with a simple
command. Enjoy the *twelve-factor app* workflow while keeping total
control.

The `opdemand/deis <https://github.com/opdemand/deis>`__ project
contains a command-line interface to the **Deis** system. It's all you
need to work with an existing **Deis** controller. To set up your own
private application platform, the
`deis-controller <https://github.com/opdemand/deis-controller>`__ and
`deis-chef <https://github.com/opdemand/deis-chef>`__ projects are also
required.

Getting Started
---------------

Clone the git repository at `<https://github.com/opdemand/deis.git>`_:

::

    $ git clone https://github.com/opdemand/deis.git
    Cloning into 'deis'...
    ...
    Resolving deltas: 100%, done.
    $ cd deis


Use an isolated python environment:

::

    $ virtualenv venv
    New python executable in venv/bin/python
    ...
    Installing pip................done.
    $ source venv/bin/activate
    (venv) $


Next Steps
----------

License
-------

**Deis** is open source software under the Apache 2.0 license. Please
see the **LICENSE** file in the root directory for details.

Credits
-------

**Deis** rests on the shoulders of leading open source technologies:

-  Docker
-  Chef
-  Django
-  Heroku buildpacks
-  Gitosis

`OpDemand <http://www.opdemand.com/>`__ sponsors and maintains the
**Deis** project.
