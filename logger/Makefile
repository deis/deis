build:
	vagrant ssh -c 'cd share/logger && sudo docker build -t deis/logger .'

install:
	vagrant ssh -c 'sudo systemctl enable /home/core/share/logger/systemd/*'

uninstall: stop
	vagrant ssh -c 'sudo systemctl disable /home/core/share/logger/systemd/*'

start:
	vagrant ssh -c 'sudo systemctl start deis-logger.service'

stop:
	vagrant ssh -c 'sudo systemctl stop deis-logger.service'

restart:
	vagrant ssh -c 'sudo systemctl restart deis-logger.service'

logs:
	vagrant ssh -c 'sudo journalctl -f -u deis-logger.service'

run: install restart logs

clean: uninstall
	vagrant ssh -c 'sudo docker rm -f deis-logger'

full-clean: clean
	vagrant ssh -c 'sudo docker rmi deis/logger'
