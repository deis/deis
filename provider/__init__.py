import importlib


def import_provider_module(provider_type):
    """
    Return the module for a provider.
    """
    try:
        tasks = importlib.import_module('provider.' + provider_type)
    except ImportError as e:
        raise e
    return tasks
