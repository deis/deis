# security keys and auth tokens
SECRET_KEY = '{{ getv "/deis/controller/secretKey" }}'
BUILDER_KEY = '{{ getv "/deis/controller/builderKey" }}'

# scheduler settings
SCHEDULER_MODULE = 'scheduler.fleet'
SCHEDULER_TARGET = '/var/run/fleet.sock'
try:
    SCHEDULER_OPTIONS = dict('{{ if exists "/deis/controller/schedulerOptions" }}{{ getv "/deis/controller/schedulerOptions" }}{{ else }}{}{{ end }}')
except:
    SCHEDULER_OPTIONS = {}

# base64-encoded SSH private key to facilitate current version of "deis run"
SSH_PRIVATE_KEY = """{{ if exists "/deis/platform/sshPrivateKey" }}{{ getv "/deis/platform/sshPrivateKey" }}{{ else }}""{{end}}"""

# platform domain must be provided
DEIS_DOMAIN = '{{ getv "/deis/platform/domain" }}'

ENABLE_PLACEMENT_OPTIONS = """{{ if exists "/deis/platform/enablePlacementOptions" }}{{ getv "/deis/platform/enablePlacementOptions" }}{{ else }}false{{end}}"""

# use the private registry module
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

LOGGER_HOST = '{{ getv "/deis/logs/host"}}'

{{ if exists "/deis/controller/registrationMode" }}
REGISTRATION_MODE = '{{ getv "/deis/controller/registrationMode" }}'
{{ end }}

{{ if exists "/deis/controller/webEnabled" }}
WEB_ENABLED = bool({{ getv "/deis/controller/webEnabled" }})
{{ end }}
UNIT_HOSTNAME = '{{ if exists "/deis/controller/unitHostname" }}{{ getv "/deis/controller/unitHostname" }}{{ else }}default{{ end }}'

{{ if exists "/deis/controller/subdomain" }}
DEIS_RESERVED_NAMES = ['{{ getv "/deis/controller/subdomain" }}']
{{ end }}

# AUTH
# LDAP
{{ if exists "/deis/controller/auth/ldap/endpoint" }}
LDAP_ENDPOINT = '{{ if exists "/deis/controller/auth/ldap/endpoint" }}{{ getv "/deis/controller/auth/ldap/endpoint"}}{{ else }} {{ end }}'
BIND_DN = '{{ if exists "/deis/controller/auth/ldap/bind/dn" }}{{ getv "/deis/controller/auth/ldap/bind/dn"}}{{ else }} {{ end }}'
BIND_PASSWORD = '{{ if exists "/deis/controller/auth/ldap/bind/password" }}{{ getv "/deis/controller/auth/ldap/bind/password"}}{{ else }} {{ end }}'
USER_BASEDN = '{{ if exists "/deis/controller/auth/ldap/user/basedn" }}{{ getv "/deis/controller/auth/ldap/user/basedn"}}{{ else }} {{ end }}'
USER_FILTER = '{{ if exists "/deis/controller/auth/ldap/user/filter" }}{{ getv "/deis/controller/auth/ldap/user/filter"}}{{ else }} {{ end }}'
GROUP_BASEDN = '{{ if exists "/deis/controller/auth/ldap/group/basedn" }}{{ getv "/deis/controller/auth/ldap/group/basedn"}}{{ else }} {{ end }}'
GROUP_FILTER = '{{ if exists "/deis/controller/auth/ldap/group/filter" }}{{ getv "/deis/controller/auth/ldap/group/filter"}}{{ else }} {{ end }}'
GROUP_TYPE = '{{ if exists "/deis/controller/auth/ldap/group/type" }}{{ getv "/deis/controller/auth/ldap/group/type"}}{{ else }} {{ end }}'
{{ end }}
