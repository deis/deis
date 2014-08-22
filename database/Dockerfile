FROM deis/base:latest
MAINTAINER OpDemand <info@opdemand.com>

# install 9.3 from postgresql.org repository
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ precise-pgdg main" > /etc/apt/sources.list.d/pgdg.list
RUN wget --quiet --no-check-certificate -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN apt-get update && apt-get install -yq postgresql-9.3 && /etc/init.d/postgresql stop

# define the execution environment
VOLUME ["/var/lib/postgresql"]
WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 5432
ADD . /app
