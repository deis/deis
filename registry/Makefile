build:
	vagrant ssh -c 'cd share/registry && sudo docker build -t deis/registry .'

install:
	vagrant ssh -c 'sudo systemctl enable /home/core/share/registry/systemd/*'

uninstall: stop
	vagrant ssh -c 'sudo systemctl disable /home/core/share/registry/systemd/*'

start:
	vagrant ssh -c 'sudo systemctl start deis-registry.service'

stop:
	-vagrant ssh -c 'sudo systemctl stop deis-registry.service'

restart:
	vagrant ssh -c 'sudo systemctl restart deis-registry.service'

logs:
	vagrant ssh -c 'sudo journalctl -f -u deis-registry.service'

run: install restart logs

clean: uninstall
	vagrant ssh -c 'sudo docker rm -f deis-registry'

full-clean: clean
	vagrant ssh -c 'sudo docker rmi deis/registry'
