FROM ubuntu:14.04

ENV DEBIAN_FRONTEND noninteractive

# install common packages
RUN apt-get update && apt-get install -y curl

# install etcdctl
RUN curl -sSL -o /usr/local/bin/etcdctl https://s3-us-west-2.amazonaws.com/opdemand/etcdctl-v0.4.6 \
    && chmod +x /usr/local/bin/etcdctl

# install confd
RUN curl -sSL -o /usr/local/bin/confd https://s3-us-west-2.amazonaws.com/opdemand/confd-v0.5.0-json \
    && chmod +x /usr/local/bin/confd

# install required system packages
# HACK: install git so we can install bacongobbler's fork of django-fsm
# install openssh-client for temporary fleetctl wrapper
RUN apt-get update && \
    apt-get install -yq python-dev libpq-dev libyaml-dev git openssh-client

# install pip
RUN curl -sSL https://raw.githubusercontent.com/pypa/pip/1.5.6/contrib/get-pip.py | python -

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

ADD . /app

# Create static resources
RUN /app/manage.py collectstatic --settings=deis.settings --noinput
