#!/bin/bash -e

VP_HOST="$(hostname | sed -e 's/^ *//g' -e 's/ *$//g').local"
VP_USER=$(whoami | sed -e 's/^ *//g' -e 's/ *$//g')

vagrant ssh -c "sudo VP_HOST=$VP_HOST VP_USER=$VP_USER /vagrant/contrib/vagrant/util/_controller-reset-db.sh"

echo "Logging in user 'devuser' with password 'devpass'..."
deis login "deis-controller.local:8000" --username=devuser --password=devpass
