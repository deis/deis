FROM deis/base:latest
MAINTAINER OpDemand <info@opdemand.com>

# install required system packages
# HACK: install git so we can install bacongobbler's fork of django-fsm
# install openssh-client for temporary fleetctl wrapper
# install recent pip
# add a deis user that has passwordless sudo (for now)
# create a /app directory for storing application data
# create directory for confd templates
# create directory for logs
RUN apt-get update && \
    apt-get install -yq python-dev libpq-dev libyaml-dev git openssh-client && \
    wget -qO- https://raw.githubusercontent.com/pypa/pip/1.5.5/contrib/get-pip.py | python - && \
    useradd deis --groups sudo --home-dir /app --shell /bin/bash && \
    sed -i -e 's/%sudo\tALL=(ALL:ALL) ALL/%sudo\tALL=(ALL:ALL) NOPASSWD:ALL/' /etc/sudoers && \
    mkdir -p /app && chown -R deis:deis /app && \
    mkdir -p /templates && chown -R deis:deis /templates && \
    mkdir -p /var/log/deis && chown -R deis:deis /var/log/deis

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
