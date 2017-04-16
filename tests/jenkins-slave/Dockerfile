FROM ubuntu-debootstrap:14.04

ENV DEBIAN_FRONTEND noninteractive

# configure locale
RUN locale-gen en_US.UTF-8 && \
    update-locale LANG=en_US.UTF-8 && \
    update-locale LC_ALL=en_US.UTF-8

# install common packages
RUN apt-get update && apt-get install -y curl

# Add docker repo
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 36A1D7869245C8950F966E92D8576A8BA88D21E9
RUN echo "deb http://get.docker.io/ubuntu docker main" > /etc/apt/sources.list.d/docker.list

# Add virtualbox repo
RUN curl -sSL http://download.virtualbox.org/virtualbox/debian/oracle_vbox.asc | apt-key add -
RUN echo "deb http://download.virtualbox.org/virtualbox/debian trusty contrib" > /etc/apt/sources.list.d/virtualbox.list

RUN apt-get update && \
    apt-get install -yq build-essential \
        mercurial git subversion wget bzr \
        python-dev libpq-dev libyaml-dev postgresql postgresql-client \
        lxc-docker-1.3.0 openjdk-7-jre-headless virtualbox-4.3

RUN wget https://dl.bintray.com/mitchellh/vagrant/vagrant_1.6.5_x86_64.deb && \
    dpkg -i vagrant_1.6.5_x86_64.deb && rm vagrant_1.6.5_x86_64.deb

# install go
RUN curl -sSL https://storage.googleapis.com/golang/go1.3.1.linux-amd64.tar.gz | tar -C /usr/local -xz
ENV PATH /usr/local/go/bin:/usr/local/bin:/usr/local/sbin:/usr/bin:/usr/sbin:/bin:/sbin
ENV GOPATH /go
ENV GOROOT /usr/local/go

# install pip and virtualenv
RUN curl -sSL https://raw.githubusercontent.com/pypa/pip/1.5.6/contrib/get-pip.py | python - && \
    pip install virtualenv

# set up PostgreSQL requirements for controller unit tests
USER postgres
RUN service postgresql start && createuser --createdb jenkins
# sudo -u postgres psql
# postgres=# create database deis owner jenkins

# create jenkins user and install node bootstrap script
USER root
RUN useradd -G docker,vboxusers -s /bin/bash -m jenkins

USER jenkins
WORKDIR /home/jenkins
RUN wget -N http://ci.deis.io/jnlpJars/slave.jar
CMD ["/bin/bash", "-c", "java -jar slave.jar -jnlpUrl http://ci.deis.io/computer/${NODE_NAME?}/slave-agent.jnlp -secret ${NODE_SECRET?}"]
