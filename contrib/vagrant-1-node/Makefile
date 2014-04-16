#
# Deis Makefile
#

# ordered list of deis components
COMPONENTS=registry logger database cache controller builder router

all: build run

test_client:
	python -m unittest ../../discover client.tests

pull:
	vagrant ssh -c 'for c in $(COMPONENTS); do docker pull deis/$$c; done'

build:
	vagrant ssh -c 'cd share && for c in $(COMPONENTS); do cd $$c && docker build -t deis/$$c . && cd ..; done'

install:
	vagrant ssh -c 'cd share && for c in $(COMPONENTS); do cd $$c && sudo systemctl enable $$(pwd)/systemd/* && cd ..; done'

uninstall: stop
	vagrant ssh -c 'cd share && for c in $(COMPONENTS); do cd $$c && sudo systemctl disable $$(pwd)/systemd/* && cd ..; done'

start:
	vagrant ssh -c 'cd share && for c in $(COMPONENTS); do cd $$c/systemd && sudo systemctl start * && cd ../..; done'

stop:
	vagrant ssh -c 'cd share && for c in $(COMPONENTS); do cd $$c/systemd && sudo systemctl stop * && cd ../..; done'

restart:
	vagrant ssh -c 'cd share && for c in $(COMPONENTS); do cd $$c/systemd && sudo systemctl restart * && cd ../..; done'

logs:
	vagrant ssh -c 'journalctl -f -u deis-*'

run: install restart logs

clean: uninstall
	vagrant ssh -c 'cd share && for c in $(COMPONENTS); do docker rm -f deis-$$c; done'

full-clean: clean
	vagrant ssh -c 'cd share && for c in $(COMPONENTS); do docker rmi deis-$$c; done'
