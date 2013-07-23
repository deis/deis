all:
	python manage.py runserver

db:
	python manage.py syncdb --noinput

test:
	python manage.py test api web

task:
	python manage.py test celerytasks

pep8:
	pep8 api celerytasks deis web

pyflakes:
	pyflakes api celerytasks deis web
