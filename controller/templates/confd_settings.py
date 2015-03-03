# security keys and auth tokens
SECRET_KEY = '{{ getv "/deis/controller/secretKey" }}'
BUILDER_KEY = '{{ getv "/deis/controller/builderKey" }}'

# scheduler settings
SCHEDULER_MODULE = 'scheduler.{{ or (.deis_controller_schedulerModule) "fleet" }}'
SCHEDULER_TARGET = '{{ or (.deis_controller_schedulerTarget) "/var/run/fleet.sock" }}'
try:
    SCHEDULER_OPTIONS = dict('{{ or (.deis_controller_schedulerOptions) "{}" }}')
except:
    SCHEDULER_OPTIONS = {}

# base64-encoded SSH private key to facilitate current version of "deis run"
SSH_PRIVATE_KEY = """{{ or (.deis_platform_sshPrivateKey) "" }}"""

# platform domain must be provided
DEIS_DOMAIN = '{{ getv "/deis/platform/domain" }}'

# use the private registry module
REGISTRY_MODULE = 'registry.private'
REGISTRY_URL = '{{ getv "/deis/registry/protocol" }}://{{ getv "/deis/registry/host" }}:{{ getv "/deis/registry/port" }}'  # noqa
REGISTRY_HOST = '{{ getv "/deis/registry/host" }}'
REGISTRY_PORT = '{{ getv "/deis/registry/port" }}'

# default to sqlite3, but allow postgresql config through envvars
DATABASES = {
    'default': {
        'ENGINE': 'django.db.backends.{{ getv "/deis/database/engine" }}',
        'NAME': '{{ getv "/deis/database/name" }}',
        'USER': '{{ getv "/deis/database/user" }}',
        'PASSWORD': '{{ getv "/deis/database/password" }}',
        'HOST': '{{ getv "/deis/database/host" }}',
        'PORT': '{{ getv "/deis/database/port" }}',
    }
}

# move log directory out of /app/deis
DEIS_LOG_DIR = '/data/logs'

{{ if exists "/deis/controller/registrationEnabled" }}
REGISTRATION_ENABLED = bool({{ getv "/deis/controller/registrationEnabled" }})
{{ end }}

{{ if exists "/deis/controller/webEnabled" }}
WEB_ENABLED = bool({{ getv "/deis/controller/webEnabled" }})
{{ end }}

UNIT_HOSTNAME = '{{ or (.deis_controller_unitHostname) "default" }}'
