"""
Mock service commands for testing
"""

DASHBOARD = '/services/mock/dashboard'
DOCS = '/services/mock/docs'
DESCRIPTION = 'Mock service for testing purposes'
PLANS = ('free', 'paid')


def build_service(service):
    return "mock://asfdfsdfg:hfghgdfg@localhost:1337/{}".format(service['name'])


def destroy_service(service):
    pass


def update_service(old_service, new_service):
    return "mock://asfdfsdfg:hfghgdfg@localhost:1337/{}".format(new_service['name'])
