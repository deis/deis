build:
	vagrant ssh -c 'cd share/database && sudo docker build -t deis/database .'

install:
	vagrant ssh -c 'sudo systemctl enable /home/core/share/database/systemd/*'

uninstall: stop
	vagrant ssh -c 'sudo systemctl disable /home/core/share/database/systemd/*'

start:
	vagrant ssh -c 'sudo systemctl start deis-database.service'

stop:
	vagrant ssh -c 'sudo systemctl stop deis-database.service'

restart:
	vagrant ssh -c 'sudo systemctl restart deis-database.service'

logs:
	vagrant ssh -c 'sudo journalctl -f -u deis-database.service'

run: install restart logs

clean: uninstall
	vagrant ssh -c 'sudo docker rm -f deis-database'

full-clean: clean
	vagrant ssh -c 'sudo docker rmi deis/database'
