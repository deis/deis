build:
	docker build -t deis/controller .

run:
	docker run -p $${PORT:-8000}:$${PORT:-8000} -e ETCD=$${ETCD:-127.0.0.1:4001} -name deis-controller deis/controller
	exit 0

shell:
	docker run -t -i -e ETCD=$${ETCD:-127.0.0.1:4001} deis/controller /bin/bash

clean:
	-docker rmi deis/controller

test:
	python manage.py test --noinput api web

runserver:
	python manage.py runserver

db:
	python manage.py syncdb --migrate --noinput

coverage:
	coverage run manage.py test --noinput api web
	coverage html

flake8:
	flake8
