.PHONY: all test logs

all: build run

build:
	vagrant ssh -c 'cd share/controller && sudo docker build -t deisreleases/controller-v0.8.0 .'

install:
	vagrant ssh -c 'sudo systemctl enable /home/core/share/controller/systemd/*'

uninstall: stop
	vagrant ssh -c 'sudo systemctl disable /home/core/share/controller/systemd/*'

start:
	vagrant ssh -c 'sudo systemctl start deis-controller.service'

stop:
	vagrant ssh -c 'sudo systemctl stop deis-controller.service'

restart:
	vagrant ssh -c 'sudo systemctl restart deis-controller.service'

logs:
	vagrant ssh -c 'sudo journalctl -f -u deis-controller.service'

run: install restart logs

clean: uninstall
	vagrant ssh -c 'sudo docker rm -f deis-controller'

full-clean: clean
	vagrant ssh -c 'sudo docker rmi deisreleases/controller-v0.8.0'

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
