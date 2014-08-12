COMPONENTS=builder cache controller database logger registry router

build:
	go build .

package:
	rm -f package
	docker build -t deis/deisctl .
	mkdir -p package
	-docker cp `docker run -d deis/deisctl`:/tmp/deisctl.tar.gz package/
	mv package/deisctl.tar.gz package/deisctl-v`cat deis-version`.tar.gz
