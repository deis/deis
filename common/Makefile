include ../includes.mk

build:
	$(call rsync_all)
	$(call ssh_all,'cd share/common && sudo docker build -t deis/common:$(VERSION) .')

full-clean:
	$(call ssh_all,'cd share/common && docker rmi deis/common:$(VERSION)')
