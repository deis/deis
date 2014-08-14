#!/usr/bin/env python

"""Install the Deis command-line client."""


try:
    from setuptools import setup
    USE_SETUPTOOLS = True
except ImportError:
    from distutils.core import setup
    USE_SETUPTOOLS = False

try:
    LONG_DESCRIPTION = open('README.rst').read()
except IOError:
    LONG_DESCRIPTION = 'Deis command-line client'

try:
    APACHE_LICENSE = open('LICENSE').read()
except IOError:
    APACHE_LICENSE = 'See http://www.apache.org/licenses/LICENSE-2.0'

KWARGS = {}
if USE_SETUPTOOLS:
    KWARGS = {'entry_points': {'console_scripts': ['deis = deis:main']}}
else:
    KWARGS = {'scripts': ['deis']}


setup(name='deis',
      version='0.11.0',
      license=APACHE_LICENSE,
      description='Command-line Client for Deis, the open PaaS',
      author='OpDemand',
      author_email='info@opdemand.com',
      url='https://github.com/deis/deis',
      keywords=[
          'opdemand', 'deis', 'paas', 'cloud', 'chef', 'docker', 'heroku',
          'aws', 'ec2', 'rackspace', 'digitalocean', 'gce'
      ],
      classifiers=[
          'Development Status :: 4 - Beta',
          'Environment :: Console',
          'Intended Audience :: Developers',
          'Intended Audience :: Information Technology',
          'Intended Audience :: System Administrators',
          'License :: OSI Approved :: Apache Software License',
          'Operating System :: OS Independent',
          'Programming Language :: Python',
          'Programming Language :: Python :: 2.7',
          'Topic :: Internet',
          'Topic :: System :: Systems Administration',
      ],
      py_modules=['deis'],
      data_files=[
          ('.', ['README.rst']),
      ],
      long_description=LONG_DESCRIPTION,
      install_requires=[
          'docopt==0.6.2', 'python-dateutil==2.2',
          'PyYAML==3.11', 'requests==2.3.0'
      ],
      zip_safe=True,
      **KWARGS)
