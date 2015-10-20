"""
Django settings for the Deis project.
"""

from __future__ import unicode_literals
import os.path
import random
import semantic_version as semver
import string
import sys
import tempfile
import ldap

from django_auth_ldap.config import LDAPSearch, GroupOfNamesType


PROJECT_ROOT = os.path.normpath(os.path.join(os.path.dirname(__file__), '..'))

DEBUG = False
TEMPLATE_DEBUG = DEBUG

ADMINS = (
    # ('Your Name', 'your_email@example.com'),
)

MANAGERS = ADMINS

CONN_MAX_AGE = 60 * 3

# SECURITY: change this to allowed fqdn's to prevent host poisioning attacks
# https://docs.djangoproject.com/en/1.6/ref/settings/#allowed-hosts
ALLOWED_HOSTS = ['*']

# Local time zone for this installation. Choices can be found here:
# http://en.wikipedia.org/wiki/List_of_tz_zones_by_name
# although not all choices may be available on all operating systems.
# In a Windows environment this must be set to your system time zone.
TIME_ZONE = 'UTC'

# Language code for this installation. All choices can be found here:
# http://www.i18nguy.com/unicode/language-identifiers.html
LANGUAGE_CODE = 'en-us'

SITE_ID = 1

# If you set this to False, Django will make some optimizations so as not
# to load the internationalization machinery.
USE_I18N = True

# If you set this to False, Django will not format dates, numbers and
# calendars according to the current locale.
USE_L10N = True

# If you set this to False, Django will not use timezone-aware datetimes.
USE_TZ = True

# Absolute filesystem path to the directory that will hold user-uploaded files.
# Example: "/var/www/example.com/media/"
MEDIA_ROOT = ''

# URL that handles the media served from MEDIA_ROOT. Make sure to use a
# trailing slash.
# Examples: "http://example.com/media/", "http://media.example.com/"
MEDIA_URL = ''

# Absolute path to the directory static files should be collected to.
# Don't put anything in this directory yourself; store your static files
# in apps' "static/" subdirectories and in STATICFILES_DIRS.
# Example: "/var/www/example.com/static/"
STATIC_ROOT = os.path.abspath(os.path.join(__file__, '..', '..', 'static'))

# URL prefix for static files.
# Example: "http://example.com/static/", "http://static.example.com/"
STATIC_URL = '/static/'

# Additional locations of static files
STATICFILES_DIRS = (
    # Put strings here, like "/home/html/static" or "C:/www/django/static".
    # Always use forward slashes, even on Windows.
    # Don't forget to use absolute paths, not relative paths.
)

# List of finder classes that know how to find static files in
# various locations.
STATICFILES_FINDERS = (
    'django.contrib.staticfiles.finders.FileSystemFinder',
    'django.contrib.staticfiles.finders.AppDirectoriesFinder',
)

# Make this unique, and don't share it with anybody.
SECRET_KEY = None  # @UnusedVariable

# List of callables that know how to import templates from various sources.
TEMPLATE_LOADERS = (
    'django.template.loaders.filesystem.Loader',
    'django.template.loaders.app_directories.Loader',
)

TEMPLATE_CONTEXT_PROCESSORS = (
    "django.contrib.auth.context_processors.auth",
    "django.core.context_processors.debug",
    "django.core.context_processors.i18n",
    "django.core.context_processors.media",
    "django.core.context_processors.request",
    "django.core.context_processors.static",
    "django.core.context_processors.tz",
    "django.contrib.messages.context_processors.messages",
    "deis.context_processors.site",
)

MIDDLEWARE_CLASSES = (
    'corsheaders.middleware.CorsMiddleware',
    'django.middleware.common.CommonMiddleware',
    'django.contrib.sessions.middleware.SessionMiddleware',
    'django.contrib.auth.middleware.AuthenticationMiddleware',
    'django.contrib.messages.middleware.MessageMiddleware',
    'api.middleware.APIVersionMiddleware',
    'deis.middleware.PlatformVersionMiddleware',
    # Uncomment the next line for simple clickjacking protection:
    # 'django.middleware.clickjacking.XFrameOptionsMiddleware',
)

ROOT_URLCONF = 'deis.urls'

# Python dotted path to the WSGI application used by Django's runserver.
WSGI_APPLICATION = 'deis.wsgi.application'

TEMPLATE_DIRS = (
    # Put strings here, like "/home/html/django_templates"
    # or "C:/www/django/templates".
    # Always use forward slashes, even on Windows.
    # Don't forget to use absolute paths, not relative paths.
    PROJECT_ROOT + '/web/templates',
)

INSTALLED_APPS = (
    'django.contrib.admin',
    'django.contrib.auth',
    'django.contrib.contenttypes',
    'django.contrib.humanize',
    'django.contrib.messages',
    'django.contrib.sessions',
    'django.contrib.sites',
    'django.contrib.staticfiles',
    # Third-party apps
    'django_auth_ldap',
    'guardian',
    'json_field',
    'gunicorn',
    'rest_framework',
    'rest_framework.authtoken',
    'south',
    'corsheaders',
    # Deis apps
    'api',
    'registry',
    'web',
)

AUTHENTICATION_BACKENDS = (
    "django_auth_ldap.backend.LDAPBackend",
    "django.contrib.auth.backends.ModelBackend",
    "guardian.backends.ObjectPermissionBackend",
)

ANONYMOUS_USER_ID = -1
LOGIN_URL = '/v1/auth/login/'
LOGIN_REDIRECT_URL = '/'

SOUTH_TESTS_MIGRATE = False

CORS_ORIGIN_ALLOW_ALL = True

CORS_ALLOW_HEADERS = (
    'content-type',
    'accept',
    'origin',
    'Authorization',
    'Host',
)

CORS_EXPOSE_HEADERS = (
    'X_DEIS_API_VERSION',  # DEPRECATED
    'X_DEIS_PLATFORM_VERSION',  # DEPRECATED
    'X-Deis-Release',  # DEPRECATED
    'DEIS_API_VERSION',
    'DEIS_PLATFORM_VERSION',
    'Deis-Release',
)

REST_FRAMEWORK = {
    'DEFAULT_MODEL_SERIALIZER_CLASS':
    'rest_framework.serializers.ModelSerializer',
    'DEFAULT_PERMISSION_CLASSES': (
        'rest_framework.permissions.IsAuthenticated',
    ),
    'DEFAULT_AUTHENTICATION_CLASSES': (
        'rest_framework.authentication.TokenAuthentication',
    ),
    'DEFAULT_RENDERER_CLASSES': (
        'rest_framework.renderers.JSONRenderer',
    ),
    'PAGINATE_BY': 100,
    'PAGINATE_BY_PARAM': 'page_size',
    'TEST_REQUEST_DEFAULT_FORMAT': 'json',
}

# URLs that end with slashes are ugly
APPEND_SLASH = False

# Determine where to send syslog messages
if os.path.exists('/dev/log'):           # Linux rsyslog
    SYSLOG_ADDRESS = '/dev/log'
elif os.path.exists('/var/log/syslog'):  # Mac OS X syslog
    SYSLOG_ADDRESS = '/var/log/syslog'
else:                                    # default SysLogHandler address
    SYSLOG_ADDRESS = ('localhost', 514)

# A sample logging configuration. The only tangible logging
# performed by this configuration is to send an email to
# the site admins on every HTTP 500 error when DEBUG=False.
# See http://docs.djangoproject.com/en/dev/topics/logging for
# more details on how to customize your logging configuration.
LOGGING = {
    'version': 1,
    'disable_existing_loggers': False,
    'formatters': {
        'verbose': {
            'format': '%(levelname)s %(asctime)s %(module)s %(process)d %(thread)d %(message)s'
        },
        'simple': {
            'format': '%(levelname)s %(message)s'
        },
    },
    'filters': {
        'require_debug_false': {
            '()': 'django.utils.log.RequireDebugFalse'
        }
    },
    'handlers': {
        'null': {
            'level': 'DEBUG',
            'class': 'logging.NullHandler',
        },
        'console': {
            'level': 'DEBUG',
            'class': 'logging.StreamHandler',
            'formatter': 'simple'
        },
        'mail_admins': {
            'level': 'ERROR',
            'filters': ['require_debug_false'],
            'class': 'django.utils.log.AdminEmailHandler'
        },
        'rsyslog': {
            'class': 'logging.handlers.SysLogHandler',
            'address': SYSLOG_ADDRESS,
            'facility': 'local0',
        },
    },
    'loggers': {
        'django': {
            'handlers': ['null'],
            'level': 'INFO',
            'propagate': True,
        },
        'django.request': {
            'handlers': ['console', 'mail_admins'],
            'level': 'WARNING',
            'propagate': True,
        },
        'api': {
            'handlers': ['console', 'mail_admins', 'rsyslog'],
            'level': 'INFO',
            'propagate': True,
        },
        'registry': {
            'handlers': ['console', 'mail_admins', 'rsyslog'],
            'level': 'INFO',
            'propagate': True,
        },
    }
}
TEST_RUNNER = 'api.tests.SilentDjangoTestSuiteRunner'

# etcd settings
ETCD_HOST, ETCD_PORT = os.environ.get('ETCD', '127.0.0.1:4001').split(',')[0].split(':')

# default deis settings
LOG_LINES = 1000
TEMPDIR = tempfile.mkdtemp(prefix='deis')
DEIS_DOMAIN = 'deisapp.local'

# standard datetime format used for logging, model timestamps, etc.
DEIS_DATETIME_FORMAT = '%Y-%m-%dT%H:%M:%S%Z'

# names which apps cannot reserve for routing
DEIS_RESERVED_NAMES = ['deis']

# default scheduler settings
SCHEDULER_MODULE = 'scheduler.mock'
SCHEDULER_TARGET = ''  # path to scheduler endpoint (e.g. /var/run/fleet.sock)
SCHEDULER_AUTH = ''
SCHEDULER_OPTIONS = {}

# security keys and auth tokens
SSH_PRIVATE_KEY = ''  # used for SSH connections to facilitate "deis run"
SECRET_KEY = os.environ.get('DEIS_SECRET_KEY', 'CHANGEME_sapm$s%upvsw5l_zuy_&29rkywd^78ff(qi')
BUILDER_KEY = os.environ.get('DEIS_BUILDER_KEY', 'CHANGEME_sapm$s%upvsw5l_zuy_&29rkywd^78ff(qi')

# registry settings
REGISTRY_URL = 'http://localhost:5000'
REGISTRY_HOST = 'localhost'
REGISTRY_PORT = 5000

# logger settings
LOGGER_HOST = 'localhost'
LOGGER_PORT = 8088

# check if we can register users with `deis register`
REGISTRATION_ENABLED = True

# check if we should enable the web UI module
WEB_ENABLED = False

# default to sqlite3, but allow postgresql config through envvars
DATABASES = {
    'default': {
        'ENGINE': 'django.db.backends.' + os.environ.get('DATABASE_ENGINE', 'postgresql_psycopg2'),
        'NAME': os.environ.get('DATABASE_NAME', 'deis'),
        # randomize test database name so we can run multiple unit tests simultaneously
        'TEST_NAME': "unittest-{}".format(''.join(
            random.choice(string.ascii_letters + string.digits) for _ in range(8)))
    }
}

APP_URL_REGEX = '[a-z0-9-]+'

# Honor HTTPS from a trusted proxy
# see https://docs.djangoproject.com/en/1.6/ref/settings/#secure-proxy-ssl-header
SECURE_PROXY_SSL_HEADER = ('HTTP_X_FORWARDED_PROTO', 'https')

# Unit Hostname handling.
# Supports:
#  default      - Docker generated hostname
#  application  - Hostname based on application unit name (i.e. my-application.v2.web.1)
#  server       - Hostname based on CoreOS server hostname
UNIT_HOSTNAME = 'default'

# LDAP DEFAULT SETTINGS (Overrided by confd later)
LDAP_ENDPOINT = ""
BIND_DN = ""
BIND_PASSWORD = ""
USER_BASEDN = ""
USER_FILTER = ""
GROUP_BASEDN = ""
GROUP_FILTER = ""
GROUP_TYPE = ""

# Create a file named "local_settings.py" to contain sensitive settings data
# such as database configuration, admin email, or passwords and keys. It
# should also be used for any settings which differ between development
# and production.
# The local_settings.py file should *not* be checked in to version control.
try:
    from .local_settings import *  # noqa
except ImportError:
    pass

# have confd_settings within container execution override all others
# including local_settings (which may end up in the container)
if os.path.exists('/templates/confd_settings.py'):
    sys.path.append('/templates')
    from confd_settings import *  # noqa

# Disable swap when mem limits are set, unless Docker is too old
DISABLE_SWAP = '--memory-swap=-1'
try:
    version = 'unknown'
    from registry.dockerclient import DockerClient
    version = DockerClient().client.version().get('Version')
    if not semver.validate(version) or semver.Version(version) < semver.Version('1.5.0'):
        DISABLE_SWAP = ''
except:
    print("Not disabling --memory-swap for Docker version {}".format(version))

# LDAP Backend Configuration
# Should be always after the confd_settings import.
LDAP_USER_SEARCH = LDAPSearch(
    base_dn=USER_BASEDN,
    scope=ldap.SCOPE_SUBTREE,
    filterstr="(%s=%%(user)s)" % USER_FILTER
)
LDAP_GROUP_SEARCH = LDAPSearch(
    base_dn=GROUP_BASEDN,
    scope=ldap.SCOPE_SUBTREE,
    filterstr="(%s=%s)" % (GROUP_FILTER, GROUP_TYPE)
)
AUTH_LDAP_SERVER_URI = LDAP_ENDPOINT
AUTH_LDAP_BIND_DN = BIND_DN
AUTH_LDAP_BIND_PASSWORD = BIND_PASSWORD
AUTH_LDAP_USER_SEARCH = LDAP_USER_SEARCH
AUTH_LDAP_GROUP_SEARCH = LDAP_GROUP_SEARCH
AUTH_LDAP_GROUP_TYPE = GroupOfNamesType()
AUTH_LDAP_USER_ATTR_MAP = {
    "first_name": "givenName",
    "last_name": "sn",
    "email": "mail",
    "username": USER_FILTER,
}
AUTH_LDAP_GLOBAL_OPTIONS = {
    ldap.OPT_X_TLS_REQUIRE_CERT: False,
    ldap.OPT_REFERRALS: False
}
AUTH_LDAP_ALWAYS_UPDATE_USER = True
AUTH_LDAP_MIRROR_GROUPS = True
AUTH_LDAP_FIND_GROUP_PERMS = True
AUTH_LDAP_CACHE_GROUPS = False
