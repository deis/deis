build:
	docker build -t deis/discovery .
	
run:
	docker run -p 4001:4001 -p 7001:7001 -e PUBLIC_IP=10.0.1.8 -rm deis/discovery ; exit 0

shell:
	docker run -t -i -rm -e PUBLIC_IP=10.0.1.8 deis/discovery /bin/bash

clean:
	docker rmi deis/discovery
