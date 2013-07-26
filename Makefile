all:
	python manage.py runserver

db:
	python manage.py syncdb --noinput
	python manage.py migrate

test:
	python manage.py test api celerytasks web

coverage:
	coverage run manage.py test api celerytasks web
	coverage html

flake8:
	flake8
