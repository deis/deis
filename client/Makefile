
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
