runserver:
	cd controller && python manage.py runserver

db:
	cd controller && python manage.py syncdb --migrate --noinput

test:
	cd controller && python manage.py test --noinput api cm provider web

coverage:
	cd controller && coverage run manage.py test --noinput api cm provider web
	cd controller && coverage html

test_client:
	python -m unittest discover client.tests

client_binary:
	cd client && pyinstaller deis.spec

flake8:
	flake8
