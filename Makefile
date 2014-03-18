test_client:
	python -m unittest discover client.tests

client_binary:
	cd client && pyinstaller deis.spec
