#!/usr/bin/env bash

# This code is kept in a separate file from `reste-db.sh` for no other reason than to not have to deal
# with the nightmare of double escaping all these commands through `vagrant ssh -c "\\\\\\\\\AGH!"`

# NB. Command for exporting fixtures
# `./manage.py dumpdata --natural --indent=4 -e sessions -e admin -e contenttypes -e auth.Permission -e south > /app/deis/api/fixtures/dev.json`

set -e

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

if [[ -z "$VP_HOST" && -z "$VP_USER" ]]; then
	echo "VP_HOST and VP_USER must be set"
	exit 1
fi

cd /vagrant/contrib/vagrant/util

echo "Dropping and recreating Deis database..."
echo "su postgres -c 'dropdb deis && createdb --encoding=utf8 --template=template0 deis'" | ./dshell deis-database

echo "Running South migrations..."
echo '/app/manage.py syncdb --migrate --noinput' | ./dshell deis-controller

echo "Importing fixtures"
echo '/app/manage.py loaddata /app/api/fixtures/dev.json' | ./dshell deis-controller

# Most of the fixture data is generic. However the host machine's SSH credentials will change
# from developer to developer.
echo "Updating vagrant provider details"
cat <<EOF | ./dshell deis-database
su postgres
psql deis -c "UPDATE api_provider SET creds = '{\"host\": \"$VP_HOST\", \"user\": \"$VP_USER\"}' WHERE owner_id = 1 AND type = 'vagrant' "
EOF
