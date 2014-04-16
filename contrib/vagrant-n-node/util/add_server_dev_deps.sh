#!/bin/bash
cd /vagrant/contrib/vagrant/util/

# Use the `coverage' command to signal whether the container has the dev dependencies
echo "which coverage > /dev/null" | ./dshell deis-controller
if [ $? -ne 0 ]; then
	cat <<-EOF | ./dshell deis-controller
		cd /app
		pip install -r dev_requirements.txt
	EOF
else
	echo "Deis server development dependencies already installed."
fi
