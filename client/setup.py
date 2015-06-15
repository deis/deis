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


with open('requirements.txt') as f:
    required = f.read().splitlines()
    required = [r for r in required if r.strip() and not r.startswith('#')]


setup(name='deis',
      version='1.7.3',
      license=APACHE_LICENSE,
      description='Command-line Client for Deis, the open PaaS',
      author='Engine Yard',
      author_email='info@deis.io',
      url='https://github.com/deis/deis',
      keywords=[
          'opdemand', 'deis', 'paas', 'cloud', 'coreos', 'docker', 'heroku',
          'aws', 'ec2', 'rackspace', 'digitalocean', 'gce'
      ],
      classifiers=[
          'Development Status :: 5 - Production/Stable',
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
      install_requires=required,
      zip_safe=True,
      **KWARGS)
