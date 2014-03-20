runserver:
	python manage.py runserver

db:
	python manage.py syncdb --migrate --noinput

test:
	python manage.py test --noinput api cm provider web

coverage:
	coverage run manage.py test --noinput api cm provider web
	coverage html

build:
	docker build -t deis/controller .

run:
	docker run -rm -p $${PORT:-8000}:$${PORT:-8000} -e ETCD=$${ETCD:-127.0.0.1:4001} -name deis-controller deis/controller ; exit 0

shell:
	docker run -t -i -rm -e ETCD=$${ETCD:-127.0.0.1:4001} deis/controller /bin/bash

clean:
	-docker rmi deis/controller

flake8:
	flake8
