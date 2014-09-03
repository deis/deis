FROM deis/store-base:latest
MAINTAINER OpDemand <info@opdemand.com>

# install Apache and RADOS gateway
RUN wget -q -O- https://raw.github.com/ceph/ceph/master/keys/autobuild.asc | apt-key add -
RUN echo deb http://gitbuilder.ceph.com/apache2-deb-precise-x86_64-basic/ref/master precise main > /etc/apt/sources.list.d/ceph-apache.list
RUN echo deb http://gitbuilder.ceph.com/libapache-mod-fastcgi-deb-precise-x86_64-basic/ref/master precise main > /etc/apt/sources.list.d/ceph-fastcgi.list

RUN apt-get update && apt-get install -yq apache2 libapache2-mod-fastcgi radosgw radosgw-agent

# enable required Apache modules
RUN a2enmod rewrite
RUN a2enmod fastcgi

# change the default port
RUN sed -i 's/80/8888/' /etc/apache2/ports.conf

# set the ServerName
RUN echo "ServerName deis-store-gateway" > /etc/apache2/conf.d/servername

# create gateway site, and disable the default site
ADD s3gw.fcgi /var/www/s3gw.fcgi
ADD rgw.conf /etc/apache2/sites-available/rgw.conf
RUN a2ensite rgw.conf
RUN a2dissite default

WORKDIR /app
EXPOSE 8888
CMD ["/app/bin/boot"]
ADD bin/boot /app/bin/boot
