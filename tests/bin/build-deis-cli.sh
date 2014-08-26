#!/bin/sh
#
# Creates a python virtual environment and builds the `deis` client binary with it.

virtualenv --system-site-packages venv
. venv/bin/activate
pip install docopt==0.6.2 python-dateutil==2.2 PyYAML==3.11 requests==2.3.0 pyinstaller==2.1
make -C client/ client
