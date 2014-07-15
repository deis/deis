:title: Client Reference
:description: A reference guide to all Deis client commands.

.. _client_ref:

Client Reference
================

.. _deis_apps:

deis apps
---------

.. automethod:: client.deis.DeisClient.apps_create
.. automethod:: client.deis.DeisClient.apps_list
.. automethod:: client.deis.DeisClient.apps_info
.. automethod:: client.deis.DeisClient.apps_open
.. automethod:: client.deis.DeisClient.apps_logs
.. automethod:: client.deis.DeisClient.apps_run
.. automethod:: client.deis.DeisClient.apps_destroy

.. _deis_auth:

deis auth
---------

.. automethod:: client.deis.DeisClient.auth_register
.. automethod:: client.deis.DeisClient.auth_cancel
.. automethod:: client.deis.DeisClient.auth_login
.. automethod:: client.deis.DeisClient.auth_logout

.. _deis_builds:

deis builds
-----------

.. automethod:: client.deis.DeisClient.builds_list
.. automethod:: client.deis.DeisClient.builds_create

.. _deis_clusters:

deis clusters
-------------

.. automethod:: client.deis.DeisClient.clusters_create
.. automethod:: client.deis.DeisClient.clusters_list
.. automethod:: client.deis.DeisClient.clusters_update
.. automethod:: client.deis.DeisClient.clusters_info
.. automethod:: client.deis.DeisClient.clusters_destroy

.. _deis_config:

deis config
-----------

.. automethod:: client.deis.DeisClient.config_list
.. automethod:: client.deis.DeisClient.config_set
.. automethod:: client.deis.DeisClient.config_unset

.. _deis_domains:

deis domains
------------

.. automethod:: client.deis.DeisClient.domains_add
.. automethod:: client.deis.DeisClient.domains_list
.. automethod:: client.deis.DeisClient.domains_remove

.. _deis_keys:

deis keys
---------

.. automethod:: client.deis.DeisClient.keys_list
.. automethod:: client.deis.DeisClient.keys_add
.. automethod:: client.deis.DeisClient.keys_remove

.. _deis_ps:

deis ps
-------

.. automethod:: client.deis.DeisClient.ps_list
.. automethod:: client.deis.DeisClient.ps_scale

.. _deis_releases:

deis releases
-------------

.. automethod:: client.deis.DeisClient.releases_list
.. automethod:: client.deis.DeisClient.releases_info
.. automethod:: client.deis.DeisClient.releases_rollback
