build:
	vagrant ssh -c 'cd share/router && sudo docker build -t deis/router .'

install:
	vagrant ssh -c 'sudo systemctl enable /home/core/share/router/systemd/*'

uninstall: stop
	vagrant ssh -c 'sudo systemctl disable /home/core/share/router/systemd/*'

start:
	vagrant ssh -c 'sudo systemctl start deis-router.service'

stop:
	vagrant ssh -c 'sudo systemctl stop deis-router.service'

restart:
	vagrant ssh -c 'sudo systemctl restart deis-router.service'

logs:
	vagrant ssh -c 'sudo journalctl -f -u deis-router.service'

run: install restart logs

clean: uninstall
	vagrant ssh -c 'sudo docker rm -f deis-router'

full-clean: clean
	vagrant ssh -c 'sudo docker rmi deis/router'
