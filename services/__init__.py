import importlib


def import_services_module(provider_type):
    tasks = importlib.import_module('services.' + provider_type)
    return tasks
