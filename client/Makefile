
build: setup-venv
	venv/bin/pip install docopt==0.6.2 python-dateutil==2.2 requests==2.3.0 pyinstaller==2.1 termcolor==1.1.0
	venv/bin/pyinstaller deis.spec
	chmod +x dist/deis

# Install deis client locally through pip.
install:
	pip install .

# Remove the installed copy of the deis client through pip.
uninstall:
	pip uninstall -y deis

# Clean up after setup.py droppings.
clean:
	rm -rf build/ dist/ *.egg-info

client:
	pyinstaller deis.spec

setup-venv:
	@if [ ! -d venv ]; then virtualenv venv; fi
	venv/bin/pip install -q flake8==2.2.2

test-style: setup-venv
	venv/bin/flake8
