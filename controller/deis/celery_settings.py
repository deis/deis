import djcelery

BROKER_URL = 'amqp://guest:guest@localhost:5672/'
TEST_RUNNER = 'djcelery.contrib.test_runner.CeleryTestSuiteRunner'
CELERY_RESULT_BACKEND = 'amqp'

# normally False to execute tasks asyncronously
# set to True to enable blocking execution for debugging
CELERY_ALWAYS_EAGER = False
EAGER_PROPAGATES_EXCEPTION = True

# make sure we import the task modules
CELERY_IMPORTS = ("celerytasks.ec2",
                  "celerytasks.controller",
                  "celerytasks.mock")

djcelery.setup_loader()
