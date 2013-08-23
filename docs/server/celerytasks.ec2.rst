:description: Python API Reference for the Deis celerytasks.ec2 module
:keywords: deis, celerytasks.ec2, python, celery, ec2, api

===============
celerytasks.ec2
===============

.. contents::
    :local:
.. currentmodule:: celerytasks.ec2

.. automodule:: celerytasks.ec2
    :members:
    :undoc-members:

    .. autofunction:: build_layer(layer, creds, params)
    .. autofunction:: destroy_layer(layer, creds, params)
    .. autofunction:: launch_node(node_id, creds, params, init, ssh_username, ssh_private_key)
    .. autofunction:: terminate_node(node_id, creds, params, provider_id)
    .. autofunction:: converge_node(node_id, ssh_username, fqdn, ssh_private_key, command='sudo chef-client')
    .. autofunction:: run_node(node_id, ssh_username, fqdn, ssh_private_key, docker_args, command)
