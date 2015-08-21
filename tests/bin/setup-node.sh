#!/usr/bin/env bash
#
# Preps a Ubuntu 14.04 box with requirements to run as a Jenkins node to https://ci.deis.io/
# Should be run as root.

# fail on any command exiting non-zero
set -eo pipefail

apt-get install -y apt-transport-https

# install docker
apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 \
            --recv-keys 36A1D7869245C8950F966E92D8576A8BA88D21E9
echo deb https://get.docker.com/ubuntu docker main > /etc/apt/sources.list.d/docker.list
apt-get update
apt-get install -yq lxc-docker-1.5.0

# install extra extensions (AUFS)
sudo apt-get -y install "linux-image-extra-$(uname -r)"

# install java
apt-get install -yq openjdk-7-jre-headless

apt-get install -yq build-essential \
                    libgl1-mesa-glx \
                    libpython2.7 \
                    libqt4-network \
                    libqt4-opengl \
                    libqtcore4 \
                    libqtgui4 \
                    libsdl1.2debian \
                    libvpx1 \
                    libxcursor1 \
                    libxinerama1 \
                    libxmu6 \
                    psmisc

# install virtualbox
if ! virtualbox --help &> /dev/null; then
  wget -nv http://download.virtualbox.org/virtualbox/4.3.22/virtualbox-4.3_4.3.22-98236~Ubuntu~raring_amd64.deb
  dpkg -i virtualbox-4.3_4.3.22-98236~Ubuntu~raring_amd64.deb
  rm virtualbox-4.3_4.3.22-98236~Ubuntu~raring_amd64.deb
fi

# install vagrant
if ! vagrant -v &> /dev/null; then
  wget -nv https://dl.bintray.com/mitchellh/vagrant/vagrant_1.7.2_x86_64.deb
  dpkg -i vagrant_1.7.2_x86_64.deb
  rm vagrant_1.7.2_x86_64.deb
fi

# install go
wget -nv -O- https://storage.googleapis.com/golang/go1.5.linux-amd64.tar.gz | tar -C /usr/local -xz
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
echo "You must reboot for the global $PATH changes to take effect."

# install test suite requirements
apt-get install -yq curl \
                    mercurial \
                    python-dev \
                    libffi-dev \
                    libpq-dev \
                    libyaml-dev \
                    git \
                    postgresql \
                    postgresql-client \
                    libldap2-dev \
                    libsasl2-dev

curl -sSL https://raw.githubusercontent.com/pypa/pip/7.0.3/contrib/get-pip.py | python -
pip install virtualenv

# TODO: rely on virtualenvs' pip instead of system pip on slaves
pip install -r "contrib/aws/requirements.txt"

# Use cabal (Haskell installer) to build and install ShellCheck
if ! shellcheck -V &> /dev/null; then
  apt-get install -yq cabal-install
  cabal update
  pushd /tmp
  git clone --branch v0.3.8 --single-branch https://github.com/koalaman/shellcheck.git
  pushd shellcheck
  cabal install --global
  popd +2
  apt-get purge -yq cabal-install
fi

# create jenkins user and install node bootstrap script
if ! getent passwd | cut -d: -f1 | grep -q jenkins; then
  useradd -G docker,vboxusers -s /bin/bash --system -m jenkins
fi

mkdir -p /home/jenkins/bin
wget -nv -x -O /home/jenkins/bin/start-node.sh \
      https://raw.githubusercontent.com/deis/deis/master/tests/bin/start-node.sh
chmod +x /home/jenkins/bin/start-node.sh
chown -R jenkins:jenkins /home/jenkins/bin

# as the jenkins user, do "vagrant plugin install vagrant-triggers
#   if not already installed"
su - jenkins -c "vagrant plugin list | grep -q vagrant-triggers || vagrant plugin install vagrant-triggers"

/etc/init.d/postgresql start

# set up PostgreSQL role for controller unit tests
sudo -u postgres psql -c "CREATE ROLE jenkins WITH CREATEDB LOGIN;" || true
sudo -u postgres psql -c "CREATE DATABASE deis WITH OWNER jenkins;" || true
# edit postgresql.conf and change "fsync = off", then restart postgresql.

# now the jenkins user has to export some envvars to start as a node
echo "Remaining setup:"
echo "  1. Log in as the jenkins user (sudo -i -u jenkins)"
echo "  2. Visit the nodes admin interface at https://ci.deis.io/ to find the command line for this node"
echo "  3. Export the NODE_NAME and NODE_SECRET environment variables defined there to your shell"
echo "  4. Run bin/start-node.sh to connect to Jenkins and start handling jobs"
