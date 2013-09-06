
from __future__ import unicode_literals


def configure_node(node):
    config = {}
    config['config'] = config
    config['node_id'] = node.id
    config['layer_id'] = node.layer.id
    return config


def bootstrap_node(node):
    return node


def converge_node(node):
    return '', 0


def destroy_node(node):
    return node


def destroy_formation(formation):
    return formation
