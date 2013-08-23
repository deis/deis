:description: Python API Reference for the Deis celerytasks.mock module
:keywords: deis, celerytasks.mock, python, celery, api

================
celerytasks.mock
================

.. contents::
    :local:
.. currentmodule:: celerytasks.mock

.. automodule:: celerytasks.mock
    :members:
    :undoc-members:

    .. autofunction:: build_layer(layer, creds, params)
    .. autofunction:: destroy_layer(layer, creds, params)
    .. autofunction:: launch_node(node_id, creds, params, init, ssh_username, ssh_private_key)
    .. autofunction:: terminate_node(node_id, creds, params, provider_id)
    .. autofunction:: converge_node(node_id, ssh_username, fqdn, ssh_private_key)
    .. autofunction:: run_node(node_id, ssh_username, fqdn, ssh_private_key, docker_args, command)
