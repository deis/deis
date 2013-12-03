all:
	python manage.py runserver

db:
	python manage.py syncdb --noinput
	python manage.py migrate --noinput

test:
	python -Wall manage.py test --noinput api cm provider web

coverage:
	coverage run manage.py test --noinput api cm provider web
	coverage html

test_client:
	python -Wall -m unittest discover client.tests

flake8:
	flake8
