import importlib


def import_provider_module(provider_type):
    """
    Return the module for a provider.
    """
    tasks = importlib.import_module('provider.' + provider_type)
    return tasks
