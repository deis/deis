#!/bin/bash
#
# Preps a Ubuntu 14.04 box with requirements to run as a Jenkins node to https://ci.deis.io/
# Should be run as root.

# install docker 1.3.3
apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 36A1D7869245C8950F966E92D8576A8BA88D21E9
sh -c "echo deb https://get.docker.io/ubuntu docker main > /etc/apt/sources.list.d/docker.list"
apt-get update && apt-get install -yq lxc-docker-1.3.3

# install java
apt-get install -yq openjdk-7-jre-headless

# install virtualbox 4.3.14
apt-get install -yq build-essential libgl1 libgl1-mesa-glx libpython2.7 libqt4-network libqt4-opengl \
    libqtcore4 libqtgui4 libsdl1.2debian libvpx1 libxcursor1
wget http://download.virtualbox.org/virtualbox/4.3.20/virtualbox-4.3_4.3.20-96996~Ubuntu~raring_amd64.deb
dpkg -i virtualbox-4.3_4.3.20-96996~Ubuntu~raring_amd64.deb && \
    rm virtualbox-4.3_4.3.20-96996~Ubuntu~raring_amd64.deb

# install vagrant
wget https://dl.bintray.com/mitchellh/vagrant/vagrant_1.7.1_x86_64.deb
dpkg -i vagrant_1.7.1_x86_64.deb && rm vagrant_1.7.1_x86_64.deb

# install go
wget -qO- https://storage.googleapis.com/golang/go1.3.3.linux-amd64.tar.gz | tar -C /usr/local -xz
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
echo "You must reboot for the global $PATH changes to take effect."

# install test suite requirements
apt-get install -yq curl mercurial python-dev libpq-dev libyaml-dev git postgresql postgresql-client
RUN curl -sSL https://raw.githubusercontent.com/pypa/pip/1.5.6/contrib/get-pip.py | python -
pip install virtualenv

# create jenkins user and install node bootstrap script
useradd -G docker,vboxusers -s /bin/bash -m jenkins
mkdir -p /home/jenkins/bin
wget -x -O /home/jenkins/bin/start-node.sh \
    https://raw.githubusercontent.com/deis/deis/master/tests/bin/start-node.sh
chmod +x /home/jenkins/bin/start-node.sh
chown -R jenkins:jenkins /home/jenkins/bin

# set up PostgreSQL role for controller unit tests
sudo -u postgres psql -c "CREATE ROLE deis WITH CREATEDB PASSWORD 'changeme123';"

# now the jenkins user has to export some envvars to start as a node
echo "Remaining setup:"
echo "  1. Log in as the jenkins user (sudo -i -u jenkins)"
echo "  2. Visit the nodes admin interface at https://ci.deis.io/ to find the command line for this node"
echo "  3. Export the NODE_NAME and NODE_SECRET environment variables defined there to your shell"
echo "  4. Run bin/start-node.sh to connect to Jenkins and start handling jobs"
