:description: Python API Reference for the Deis celerytasks.azuresms module
:keywords: deis, celerytasks.azuresms, python, celery, azure, api

====================
celerytasks.azuresms
====================

**Note: Windows Azure cloud support is not ready for testing yet.***

.. contents::
    :local:
.. currentmodule:: celerytasks.azuresms

.. automodule:: celerytasks.azuresms
    :members:
    :undoc-members:

    .. autofunction:: launch_node(node_id, creds, params, init, ssh_username, ssh_private_key)
    .. autofunction:: terminate_node(node_id, creds, params, provider_id)
    .. autofunction:: converge_node(node_id, ssh_username, fqdn, ssh_private_key, command='sudo chef-client')
