#!/bin/sh
#
# Creates a python virtual environment and builds the `deis` client binary with it.

virtualenv --system-site-packages venv
. venv/bin/activate
pip install docopt==0.6.2 python-dateutil==2.2 PyYAML==3.11 requests==2.5.1 git+https://github.com/pyinstaller/pyinstaller@7413317 termcolor==1.1.0
make -C client/ client
