:title: Manage an Application on Deis
:description: First steps for developers using Deis to deploy and scale applications.

.. _manage_application:

Manage an Application
=====================

Deis includes many tools for managing deployed :ref:`Applications <application>`.

Scale the Application
---------------------

Applications deployed on Deis `scale out via the process model`_. Use :code:`deis scale`
to control the number of :ref:`Containers <container>` that power your app.

.. code-block:: console

    $ deis scale web=8
    Scaling processes... but first, coffee!
    done in 20s

    === peachy-waxworks Processes

    --- web: `java -cp target/classes:target/dependency/* HelloWorld`
    web.1 up 2013-12-03T00:00:25.836Z (dev-runtime-1)
    web.2 up 2013-12-03T00:30:10.934Z (dev-runtime-1)
    web.3 up 2013-12-03T00:30:10.950Z (dev-runtime-1)
    web.4 up 2013-12-03T00:30:10.963Z (dev-runtime-1)
    web.5 up 2013-12-03T00:30:10.978Z (dev-runtime-1)
    web.6 up 2013-12-03T00:30:10.993Z (dev-runtime-1)
    web.7 up 2013-12-03T00:30:11.010Z (dev-runtime-1)
    web.8 up 2013-12-03T00:30:11.027Z (dev-runtime-1)

Scaling is managed by process types like :code:`web` or :code:`worker`, which are defined
in a `Procfile`_ in the root of your application repository.

.. note::

    Dockerfile-based applications use the :code:`cmd` process type to scale the default
    container command.

Administer the Application
--------------------------

Deis applications `use one-off processes for admin tasks`_ like database migrations and
other commands that must run against the live application.

Use :code:`deis run` to execute commands on the deployed application.

.. code-block:: console

    $ deis run ls -l
    total 28
    -rw-r--r-- 1 root root  553 Dec  2 23:59 LICENSE
    -rw-r--r-- 1 root root   60 Dec  2 23:59 Procfile
    -rw-r--r-- 1 root root   33 Dec  2 23:59 README.md
    -rw-r--r-- 1 root root 1622 Dec  2 23:59 pom.xml
    drwxr-xr-x 3 root root 4096 Dec  2 23:59 src
    -rw-r--r-- 1 root root   25 Dec  2 23:59 system.properties
    drwxr-xr-x 6 root root 4096 Dec  3 00:00 target

Share the Application
---------------------

Use :code:`deis sharing:add` to allow another Deis user to collaborate on your application.

.. code-block:: console

  $ deis sharing:add otheruser
  Adding otheruser to peachy-waxworks collaborators... done

Use :code:`deis sharing` to see who an application is currently shared with, and
:code:`deis sharing:remove` to remove a collaborator.

.. note::
    Collaborators can do anything with an application that its owner can do,
    except delete the application itself.

When working with an application that has been shared with you, clone the original
repository and add Deis' git remote entry before attempting to :code:`git push` any
changes to Deis.

.. code-block:: console

  $ git clone https://github.com/deis/example-java-jetty.git
  Cloning into 'example-java-jetty'... done
  $ cd example-java-jetty
  $ git remote add -f deis ssh://git@local.deisapp.com:2222/peachy-waxworks.git
  Updating deis
  From deis-controller.local:peachy-waxworks
   * [new branch]      master     -> deis/master

Troubleshoot the Application
----------------------------

Applications deployed on Deis `treat logs as event streams`_. Deis aggregates standard
output and standard error from every :ref:`Container`, making it easy to troubleshoot
problems with your application.

Use :ref:`deis logs <deis_apps>` to view the log output from your deployed application.

.. code-block:: console

    $ deis logs | tail
    Dec  3 00:30:31 ip-10-250-15-201 peachy-waxworks[web.5]: INFO:oejsh.ContextHandler:started o.e.j.s.ServletContextHandler{/,null}
    Dec  3 00:30:31 ip-10-250-15-201 peachy-waxworks[web.8]: INFO:oejs.Server:jetty-7.6.0.v20120127
    Dec  3 00:30:31 ip-10-250-15-201 peachy-waxworks[web.5]: INFO:oejs.AbstractConnector:Started SelectChannelConnector@0.0.0.0:10005
    Dec  3 00:30:31 ip-10-250-15-201 peachy-waxworks[web.6]: INFO:oejsh.ContextHandler:started o.e.j.s.ServletContextHandler{/,null}
    Dec  3 00:30:31 ip-10-250-15-201 peachy-waxworks[web.7]: INFO:oejsh.ContextHandler:started o.e.j.s.ServletContextHandler{/,null}
    Dec  3 00:30:31 ip-10-250-15-201 peachy-waxworks[web.6]: INFO:oejs.AbstractConnector:Started SelectChannelConnector@0.0.0.0:10006
    Dec  3 00:30:31 ip-10-250-15-201 peachy-waxworks[web.8]: INFO:oejsh.ContextHandler:started o.e.j.s.ServletContextHandler{/,null}
    Dec  3 00:30:31 ip-10-250-15-201 peachy-waxworks[web.7]: INFO:oejs.AbstractConnector:Started SelectChannelConnector@0.0.0.0:10007
    Dec  3 00:30:31 ip-10-250-15-201 peachy-waxworks[web.8]: INFO:oejs.AbstractConnector:Started SelectChannelConnector@0.0.0.0:10008

.. _`store config in environment variables`: http://12factor.net/config
.. _`decoupled from the application`: http://12factor.net/backing-services
.. _`scale out via the process model`: http://12factor.net/concurrency
.. _`treat logs as event streams`: http://12factor.net/logs
.. _`use one-off processes for admin tasks`: http://12factor.net/admin-processes
.. _`Procfile`: http://ddollar.github.io/foreman/#PROCFILE
