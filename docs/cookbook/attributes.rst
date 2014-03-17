:title: Default Attributes
:description: Describes the default attributes available in the Chef cookbook
:keywords: deis, documentation, cookbook, chef

.. _attributes:

Attributes
==========

This is an exhaustive list of default attributes for the `deis cookbook`_.

Deis
----

dir
~~~

* ``default.deis.dir``
* ``type: String``

The home directory for the default deis user.


username
~~~~~~~~

* ``default.deis.username``
* ``type: String``

The username for the default deis user.

group
~~~~~

* ``default.deis.group``
* ``type: String``

The default group that the deis user is added to.

log_dir
~~~~~~~

* ``default.deis.log_dir``
* ``type: String``

The directory where all controller log files should go.

public_ip
~~~~~~~~~

* ``default.deis.public_ip``
* ``type: String``

The publicly addressable IP address to the cluster's controller. If this attribute is not
defined, it will be automatically discovered using Chef's `Ohai`_.

.. _image_timeout:

image_timeout
~~~~~~~~~~~~~

* ``default.deis.image_timeout``
* ``type: Integer``

The maximum length of time that a ``docker pull`` operation should run before it times out.

autoupgrade
~~~~~~~~~~~

* ``default.deis.autoupgrade``
* ``type: Boolean``

When set to true, the cookbook will redeploy containers when images are updated. It will
also make Docker pull the latest images from the Docker index, overwriting the older
image.

dev
---

.. _devmode:

mode
~~~~

* ``default.deis.dev.mode``
* ``type: Boolean``

When set to true, the cookbook will automatically mount ``default.deis.dev.source``'s
submodules into their respective containers. For example, ``deis-server`` will
automatically have the ``server`` project mounted for development.

source
~~~~~~

* ``default.deis.dev.source``
* ``type: String``

The absolute path to the deis source code on the server. This key must be set when
:ref:`deis.dev.mode <devmode>` is set to true.

rsyslog
-------

For more information on configuring rsyslog, see
https://github.com/opscode-cookbooks/rsyslog

server_search
~~~~~~~~~~~~~

* ``default.rsyslog.server_search``
* ``type: String``

Specifies the criteria for the server search operation. Default is ``role:loghost``.


etcd
----

repository
~~~~~~~~~~

* ``default.deis.etcd.repository``
* ``type: String``

The public repository on the Docker Index to pull.

tag
~~~

* ``default.deis.etcd.tag``
* ``type: String``

The tag that we wish to pull from the Docker Index.

image_timeout
~~~~~~~~~~~~~

* ``default.deis.etcd.image_timeout``
* ``type: Integer``

See :ref:`image_timeout`.

source
~~~~~~

* ``default.deis.etcd.source``
* ``type: String``

The source code to the etcd docker image. This variable is unused and not respected.
Use :ref:`deis.dev.mode <devmode>` instead.

container
~~~~~~~~~

* ``default.deis.etcd.container``
* ``type: String``

The canonical name given to the docker container running this application.

port
~~~~

* ``default.deis.etcd.port``
* ``type: Integer``

The port that etcd should listen for incoming requests.

peer_port
~~~~~~~~~

* ``default.deis.etcd.peer_port``
* ``type: Integer``

The port that etcd should listen for peer connections.

url
~~~

* ``default.deis.etcd.url``
* ``type: String``

The URL to a tarball release of etcd known to work with Deis.

Database
--------

repository
~~~~~~~~~~

* ``default.deis.database.repository``
* ``type: String``

The public repository on the Docker Index to pull.

tag
~~~

* ``default.deis.database.tag``
* ``type: String``

The tag that we wish to pull from the Docker Index.

image_timeout
~~~~~~~~~~~~~

* ``default.deis.database.image_timeout``
* ``type: Integer``

See :ref:`image_timeout`.

source
~~~~~~

* ``default.deis.database.source``
* ``type: String``

The source code to the database docker image. This variable is unused and not respected.
Use :ref:`deis.dev.mode <devmode>` instead.

container
~~~~~~~~~

* ``default.deis.database.container``
* ``type: String``

The canonical name given to the docker container running the database.

port
~~~~

* ``default.deis.database.port``
* ``type: Integer``

The port that the database should listen for incoming requests.

database_data
-------------

repository
~~~~~~~~~~

* ``default.deis.database_data.repository``
* ``type: String``

The public repository on the Docker Index to pull.

tag
~~~

* ``default.deis.database_data.tag``
* ``type: String``

The tag that we wish to pull from the Docker Index.

image_timeout
~~~~~~~~~~~~~

* ``default.deis.database_data.image_timeout``
* ``type: Integer``

See :ref:`image_timeout`.

container
~~~~~~~~~

* ``default.deis.database_data.container``
* ``type: String``

The canonical name given to the docker container running this application.

cache
-----

repository
~~~~~~~~~~

* ``default.deis.cache.repository``
* ``type: String``

The public repository on the Docker Index to pull.

tag
~~~

* ``default.deis.cache.tag``
* ``type: String``

The tag that we wish to pull from the Docker Index.

image_timeout
~~~~~~~~~~~~~

* ``default.deis.cache.image_timeout``
* ``type: Integer``

See :ref:`image_timeout`.

source
~~~~~~

* ``default.deis.cache.source``
* ``type: String``

The source code to the cache docker image. This variable is unused and not respected.
Use :ref:`deis.dev.mode <devmode>` instead.

container
~~~~~~~~~

* ``default.deis.cache.container``
* ``type: String``

The canonical name given to the docker container running this application.

port
~~~~

* ``default.deis.cache.port``
* ``type: Integer``

The port that the cache should listen for incoming requests.

server
------

repository
~~~~~~~~~~

* ``default.deis.server.repository``
* ``type: String``

The public repository on the Docker Index to pull.

tag
~~~

* ``default.deis.server.tag``
* ``type: String``

The tag that we wish to pull from the Docker Index.

image_timeout
~~~~~~~~~~~~~

* ``default.deis.server.image_timeout``
* ``type: Integer``

See :ref:`image_timeout`.

source
~~~~~~

* ``default.deis.server.source``
* ``type: String``

The source code to the controller docker image. This variable is unused and not respected.
Use :ref:`deis.dev.mode <devmode>` instead.

container
~~~~~~~~~

* ``default.deis.server.container``
* ``type: String``

The canonical name given to the docker container running this application.

port
~~~~

* ``default.deis.server.port``
* ``type: Integer``

The port that the server should listen for incoming requests.

worker
------

repository
~~~~~~~~~~

* ``default.deis.worker.repository``
* ``type: String``

The public repository on the Docker Index to pull.

tag
~~~

* ``default.deis.worker.tag``
* ``type: String``

The tag that we wish to pull from the Docker Index.

image_timeout
~~~~~~~~~~~~~

* ``default.deis.worker.image_timeout``
* ``type: Integer``

See :ref:`image_timeout`.

source
~~~~~~

* ``default.deis.worker.source``
* ``type: String``

The source code to the worker docker image. This variable is unused and not respected.
Use :ref:`deis.dev.mode <devmode>` instead.

container
~~~~~~~~~

* ``default.deis.worker.container``
* ``type: String``

The canonical name given to the docker container running this application.

registry
--------

repository
~~~~~~~~~~

* ``default.deis.registry.repository``
* ``type: String``

The public repository on the Docker Index to pull.

tag
~~~

* ``default.deis.registry.tag``
* ``type: String``

The tag that we wish to pull from the Docker Index.

image_timeout
~~~~~~~~~~~~~

* ``default.deis.registry.image_timeout``
* ``type: Integer``

See :ref:`image_timeout`.

source
~~~~~~

* ``default.deis.registry.source``
* ``type: String``

The source code to the registry docker image. This variable is unused and not respected.
Use :ref:`deis.dev.mode <devmode>` instead.

container
~~~~~~~~~

* ``default.deis.registry.container``
* ``type: String``

The canonical name given to the docker container running this application.

port
~~~~

* ``default.deis.registry.port``
* ``type: Integer``

The port that the server should listen for incoming requests.

settings_flavor
~~~~~~~~~~~~~~~

* ``default.deis.registry.settings_flavor``
* ``type: String``

The mode or flavor that you wish to run the docker registry under. Can be one of ``dev``,
``prod``, ``swift``, ``openstack``, or ``openstack-swift``.

S3
--

access_key
~~~~~~~~~~

* ``default.deis.registry.s3.access_key``
* ``type: String``

Your Amazon access key.

secret_key
~~~~~~~~~~

* ``default.deis.registry.s3.secret_key``
* ``type: String``

Your Amazon secret key.

bucket
~~~~~~

* ``default.deis.registry.s3.bucket``
* ``type: String``

The S3 bucket that will store your images.

encrypt
~~~~~~~

* ``default.deis.registry.s3.encrypt``
* ``type: Boolean``

If true, the container will be encrypted on the server-side by S3 and will be stored in an
encrypted form while at rest in S3.

secure
~~~~~~

* ``default.deis.registry.s3.secure``
* ``type: Boolean``

If true, all communication with S3 will be done over HTTPS instead of HTTP.

SMTP
----

host
~~~~

* ``default.deis.registry.smtp.host``
* ``type: String``

The SMTP hostname to connect to.

port
~~~~

* ``default.deis.registry.smtp.port``
* ``type: Integer``

The SMTP port to connect to.

login
~~~~~

* ``default.deis.registry.smtp.login``
* ``type: String``

The username to use when connecting to an authenticated SMTP host.

password
~~~~~~~~

* ``default.deis.registry.smtp.password``
* ``type: String``

The password to use when connecting to an authenticated SMTP host.

secure
~~~~~~

* ``default.deis.registry.smtp.secure``
* ``type: Boolean``

If set to true, the registry will use TLS to communicate with the SMTP server.

from
~~~~

* ``default.deis.registry.smtp.from``
* ``type: String``

The email address used when sending email.

to
~~

* ``default.deis.registry.smtp.to``
* ``type: String``

The email address to send exceptions to.

Swift
-----

auth_url
~~~~~~~~

* ``default.deis.registry.swift.auth_url``
* ``type: String``

The authentication URL for the keystone server in your openstack cluster.

container
~~~~~~~~~

* ``default.deis.registry.swift.container``
* ``type: String``

The Swift container to store images and repositories.

username
~~~~~~~~

* ``default.deis.registry.swift.username``
* ``type: String``

The username to authenticate against the keystone server.

password
~~~~~~~~

* ``default.deis.registry.swift.password``
* ``type: String``

The password to authenticate against the keystone server.

tenant_name
~~~~~~~~~~~

* ``default.deis.registry.swift.tenant_name``
* ``type: String``

The tenant name that your user is bound to.

region_name
~~~~~~~~~~~

* ``default.deis.registry.swift.region_name``
* ``type: String``

The region name that your service is available.

registry_data
-------------

repository
~~~~~~~~~~

* ``default.deis.registry_data.repository``
* ``type: String``

The public repository on the Docker Index to pull.

tag
~~~

* ``default.deis.registry_data.tag``
* ``type: String``

The tag that we wish to pull from the Docker Index.

image_timeout
~~~~~~~~~~~~~

* ``default.deis.registry_data.image_timeout``
* ``type: Integer``

See :ref:`image_timeout`.

source
~~~~~~

* ``default.deis.registry_data.source``
* ``type: String``

The source code to the registry_data docker image. This variable is unused and not
respected. Use :ref:`deis.dev.mode <devmode>` instead.

container
~~~~~~~~~

* ``default.deis.registry_data.container``
* ``type: String``

The canonical name given to the docker container running this application.

builder
-------

repository
~~~~~~~~~~

* ``default.deis.builder.repository``
* ``type: String``

The public repository on the Docker Index to pull.

tag
~~~

* ``default.deis.builder.tag``
* ``type: String``

The tag that we wish to pull from the Docker Index.

image_timeout
~~~~~~~~~~~~~

* ``default.deis.builder.image_timeout``
* ``type: Integer``

See :ref:`image_timeout`.

source
~~~~~~

* ``default.deis.builder.source``
* ``type: String``

The source code to the builder docker image. This variable is unused and not respected.
Use :ref:`deis.dev.mode <devmode>` instead.

container
~~~~~~~~~

* ``default.deis.builder.container``
* ``type: String``

The canonical name given to the docker container running this application.

port
~~~~

* ``default.deis.builder.port``
* ``type: Integer``

The port that the server should listen for incoming requests.

packs
~~~~~

* ``default.deis.builder.packs``
* ``type: String``

If set, the directory on the controller that the deis builtin buildpacks should be
synchronized to. This will also mount the specified directory into the ``deis-builder``
image.

logger
------

repository
~~~~~~~~~~

* ``default.deis.logger.repository``
* ``type: String``

The public repository on the Docker Index to pull.

tag
~~~

* ``default.deis.logger.tag``
* ``type: String``

The tag that we wish to pull from the Docker Index.

image_timeout
~~~~~~~~~~~~~

* ``default.deis.logger.image_timeout``
* ``type: Integer``

See :ref:`image_timeout`.

source
~~~~~~

* ``default.deis.logger.source``
* ``type: String``

The source code to the logger docker image. This variable is unused and not respected.
Use :ref:`deis.dev.mode <devmode>` instead.

container
~~~~~~~~~

* ``default.deis.logger.container``
* ``type: String``

The canonical name given to the docker container running this application.

port
~~~~

* ``default.deis.logger.port``
* ``type: Integer``

The port that the server should listen for incoming requests.

user
~~~~

* ``default.deis.logger.user``
* ``type: String``

The user that the logger image should run under. This variable is unused.

.. _`deis cookbook`: https://github.com/opdemand/deis-cookbook.git
.. _`Ohai`: http://docs.opscode.com/chef/ohai.html
