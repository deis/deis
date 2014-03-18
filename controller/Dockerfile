FROM deis/base:latest
MAINTAINER Gabriel A. Monroy <gabriel@opdemand.com>

# install required system packages
RUN apt-get update
RUN apt-get install -yq python-pip python-dev libpq-dev

# install chef
RUN apt-get install -yq ruby1.9.1 rubygems
RUN gem install --no-ri --no-rdoc chef

# install requirements before ADD to cache layer and speed build
RUN pip install boto==2.23.0 celery==3.1.8 Django==1.6.2 django-allauth==0.15.0 django-guardian==1.1.1 django-json-field==0.5.5 django-yamlfield==0.5 djangorestframework==2.3.12 dop==0.1.4 gevent==1.0 gunicorn==18.0 paramiko==1.12.1 psycopg2==2.5.2 pycrypto==2.6.1 python-etcd==0.3.0 pyrax==1.6.2 PyYAML==3.10 redis==2.8.0 static==1.0.2 South==0.8.4

# clone the project into /app
ADD . /app

# install python requirements
RUN pip install -r /app/requirements.txt

# Create static resources
RUN /app/manage.py collectstatic --settings=deis.settings --noinput

# add a deis user that has passwordless sudo (for now)
RUN useradd deis --groups sudo --home-dir /app --shell /bin/bash
RUN sed -i -e 's/%sudo\tALL=(ALL:ALL) ALL/%sudo\tALL=(ALL:ALL) NOPASSWD:ALL/' /etc/sudoers
RUN chown -R deis:deis /app

# create directory for logs
RUN mkdir -p /app/logs && chown -R deis:deis /app/logs

# define the execution environment
CMD ["/app/bin/boot"]
EXPOSE 8000
