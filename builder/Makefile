build:
	vagrant ssh -c 'cd share/builder && docker build -t deis/builder .'

install:
	vagrant ssh -c 'sudo systemctl enable /home/core/share/builder/systemd/*'

uninstall: stop
	vagrant ssh -c 'sudo systemctl disable /home/core/share/builder/systemd/*'

start:
	vagrant ssh -c 'sudo systemctl start deis-builder.service'

stop:
	vagrant ssh -c 'sudo systemctl stop deis-builder.service'

restart:
	vagrant ssh -c 'sudo systemctl restart deis-builder.service'

logs:
	vagrant ssh -c 'sudo journalctl -f -u deis-builder.service'

run: install restart logs

clean: uninstall
	vagrant ssh -c 'sudo docker rm -f deis-builder'

full-clean: clean
	vagrant ssh -c 'sudo docker rmi deis/builder'
