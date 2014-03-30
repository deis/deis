import importlib

from deis import settings

# import the registry module specified in settings
_registry_module = importlib.import_module(settings.REGISTRY_MODULE)

# expose the publish_release method publicly
publish_release = _registry_module.publish_release
