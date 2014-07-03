FROM deis/base:latest
MAINTAINER OpDemand <info@opdemand.com>

# install required system packages
RUN apt-get update && \
    apt-get install -yq python-dev libpq-dev libyaml-dev

# install recent pip
RUN wget -qO- https://raw.githubusercontent.com/pypa/pip/1.5.5/contrib/get-pip.py | python -

# HACK: install git so we can install bacongobbler's fork of django-fsm
RUN apt-get install -yq git

# install openssh-client for temporary fleetctl wrapper
RUN apt-get install -yq openssh-client

# add a deis user that has passwordless sudo (for now)
RUN useradd deis --groups sudo --home-dir /app --shell /bin/bash
RUN sed -i -e 's/%sudo\tALL=(ALL:ALL) ALL/%sudo\tALL=(ALL:ALL) NOPASSWD:ALL/' /etc/sudoers

# create a /app directory for storing application data
RUN mkdir -p /app && chown -R deis:deis /app

# create directory for confd templates
RUN mkdir -p /templates && chown -R deis:deis /templates

# create directory for logs
RUN mkdir -p /var/log/deis && chown -R deis:deis /var/log/deis

# define execution environment
CMD ["/app/bin/boot"]
EXPOSE 8000

# define work environment
WORKDIR /app

# install dependencies
ADD requirements.txt /app/requirements.txt
RUN pip install -r /app/requirements.txt

# clone the project into /app
ADD . /app

# Create static resources
RUN /app/manage.py collectstatic --settings=deis.settings --noinput
