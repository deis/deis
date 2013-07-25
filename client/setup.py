#!/usr/bin/env python

"""Install the Deis command-line client."""


import os.path
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
    KWARGS = {
        'install_requires': ['PyYAML', 'requests'],
        'entry_points': {'console_scripts': ['deis = deis.client:main']},
    }
else:
    KWARGS = {'scripts': [os.path.sep.join(['deis', 'deis'])]}


# pylint: disable=W0142
setup(name='deis',
      version='0.0.4',
      license=APACHE_LICENSE,
      description='Command-line Client for Deis',
      author='OpDemand',
      author_email='info@opdemand.com',
      url='https://github.com/opdemand/deis',
      keywords=['opdemand', 'deis', 'cloud', 'aws', 'ec2', 'heroku', 'docker'],
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
      packages=['deis'],
      data_files=[
          ('.', ['README.rst']),
          ],
      scripts=['deis/deis'],
      long_description=LONG_DESCRIPTION,
      requires=['PyYAML', 'requests'],
      zip_safe=True,
      **KWARGS)
