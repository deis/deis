FROM deis/base:latest
MAINTAINER Gabriel Monroy <gabriel@opdemand.com>

# install rsyslog from rsyslog.com
RUN echo "deb http://ppa.launchpad.net/adiscon/v7-stable/ubuntu precise main" > /etc/apt/sources.list.d/rsyslog.list
RUN echo "deb-src http://ppa.launchpad.net/adiscon/v7-stable/ubuntu precise main" >> /etc/apt/sources.list.d/rsyslog.list
RUN gpg --recv-keys --keyserver keyserver.ubuntu.com AEF0CF8E
RUN gpg --export --armor AEF0CF8E | sudo apt-key add -
RUN apt-get update
RUN apt-get install -yq --force-yes rsyslog

# create /var/log/deis for holding logs (access via bind mount)
RUN mkdir -p /var/log/deis

# install latest etcdctl including no-sync options
RUN wget -q https://s3-us-west-2.amazonaws.com/deis/etcdctl.no-sync -O /usr/local/bin/etcdctl
RUN chmod +x /usr/local/bin/etcdctl

# add the current build context to /app
ADD . /app

# prepare execution environment
WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 514
