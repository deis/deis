test_client:
	python -m unittest discover client.tests

build:
	for image in builder cache controller database discovery logger registry deis; do \
		pushd $$image; \
		docker build -t deis/$$image .; \
		popd; \
	done
