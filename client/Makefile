
build: setup-venv
	venv/bin/pip install -r requirements.txt git+https://github.com/pyinstaller/pyinstaller@32bbb954b355937ccfe377afbe56979db79a7b30
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

installer: build
	@if [ ! -d makeself ]; then git clone -b single-binary https://github.com/deis/makeself.git; fi
	PATH=./makeself:$$PATH BINARY=deis makeself.sh --bzip2 --current --nox11 dist \
		dist/deis-cli-`grep '__version__\ =\ ' deis.py | cut -d' ' -f3 | tr -d \'`-`go env GOOS`-`go env GOARCH`.run \
		"Deis CLI" "echo \
		&& echo '\033[0;36mdeis\033[0m is in the current directory. Please' \
		&& echo 'move \033[0;36mdeis\033[0m to a directory in your search PATH.' \
		&& echo \
		&& echo 'See http://docs.deis.io/ for documentation.' \
		&& echo"

setup-venv:
	@if [ ! -d venv ]; then virtualenv venv; fi

test-style: setup-venv
	venv/bin/pip install -q flake8==2.4.0
	venv/bin/flake8

test-unit:
	@echo no unit tests
