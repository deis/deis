#!/bin/bash -e

# This code is kept in a separate file from `reste-db.sh` for no other reason than to not have to deal
# with the nightmare of double escaping all these commands through `vagrant ssh -c "\\\\\\\\\AGH!"`

# NB. Command for exporting fixtures
# pg_dump --data-only --table=api_formations --table=auth_user --table=api_formation --table=api_provider --table=api_flavor deis > api/fixtures/deis_dev.sql
# And then edit the resulting SQL to remove the default anonymous user

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

if [[ -z "$VP_HOST" && -z "$VP_USER" ]]; then
	echo "VP_HOST and VP_USER must be set"
	exit 1
fi

echo "Dropping and recreating Deis database..."
su postgres -c 'dropdb deis && createdb --encoding=utf8 --template=template0 deis'
echo "Running South migrations..."
su deis -c '/opt/deis/controller/venv/bin/python /opt/deis/controller/manage.py syncdb --migrate --noinput'
echo "Updating the Django site object..."
su deis -c "psql deis -c \"UPDATE django_site SET domain = 'deis-controller.local', name = 'deis-controller.local' WHERE id = 1 \""
echo "Importing fixtures for formation and super user..."
su deis -c 'psql deis < /opt/deis/controller/api/fixtures/deis_dev.sql'
echo "Updating vagrant provider details"
su deis -c "psql deis -c \"UPDATE api_provider SET creds = '{\\\"host\\\": \\\"$VP_HOST\\\", \\\"user\\\": \\\"$VP_USER\\\"}' WHERE owner_id = 1 AND type = 'vagrant' \""
