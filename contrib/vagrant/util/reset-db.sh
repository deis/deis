#!/bin/bash -e

# NB. Command for exporting fixtures
# pg_dump --data-only --table=api_formations --table=auth_user deis > api/fixtures/deis_dev.sql
# And then edit the resulting SQL to remove the default anonymous user

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 1>&2
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
