test_client:
	python -m unittest discover client.tests

build:
	for image in builder cache controller database discovery logger registry; do \
		make -C $$image build; \
	done
