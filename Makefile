all:
	python manage.py runserver

db:
	python manage.py syncdb --noinput

test:
	python manage.py test api web

task:
	python manage.py test celerytasks

flake8:
	flake8
