# security keys and auth tokens
SECRET_KEY = '{{ .deis_controller_secretKey }}'
BUILDER_KEY = '{{ .deis_controller_builderKey }}'

# use the private registry module
REGISTRY_MODULE = 'registry.docker'
REGISTRY_URL = '{{ .deis_registry_protocol }}://{{ .deis_registry_host }}:{{ .deis_registry_port }}'  # noqa
REGISTRY_HOST = '{{ .deis_registry_host }}'
REGISTRY_PORT = '{{ .deis_registry_port }}'

# default to sqlite3, but allow postgresql config through envvars
DATABASES = {
    'default': {
        'ENGINE': 'django.db.backends.{{ .deis_database_engine }}',
        'NAME': '{{ .deis_database_name }}',
        'USER': '{{ .deis_database_user }}',
        'PASSWORD': '{{ .deis_database_password }}',
        'HOST': '{{ .deis_database_host }}',
        'PORT': '{{ .deis_database_port }}',
    }
}

# configure cache
BROKER_URL = 'redis://{{ .deis_cache_host }}:{{ .deis_cache_port }}/0'
CELERY_RESULT_BACKEND = BROKER_URL

# move log directory out of /app/deis
DEIS_LOG_DIR = '/var/log/deis'

{{ if .deis_controller_registrationEnabled }}
REGISTRATION_ENABLED = bool({{ .deis_controller_registrationEnabled }})
{{ end }}

{{ if .deis_controller_webEnabled }}
WEB_ENABLED = bool({{ .deis_controller_webEnabled }})
{{ end }}
