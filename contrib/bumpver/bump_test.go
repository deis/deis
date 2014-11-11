package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func TestHelp(t *testing.T) {
	allArgs := [][]string{{"-h"}, {"--help"}, {"help"}}
	for _, args := range allArgs {
		rc, out := commandOutput(args)
		if rc != 0 {
			t.Errorf("Return code was %d, expected 0", rc)
		}
		if !strings.Contains(out, "Updates the current semantic version") {
			t.Error(out)
		}
		if !strings.Contains(out, "Usage:") {
			t.Error(out)
		}
	}
}

func TestUsage(t *testing.T) {
	rc, out := commandOutput(nil)
	if rc != 1 {
		t.Errorf("Return code was %d, expected 1", rc)
	}
	if !strings.Contains(out, "Usage:") {
		t.Error(out)
	}
}

func TestVersion(t *testing.T) {
	args := []string{"--version"}
	rc, out := commandOutput(args)
	if rc != 0 {
		t.Errorf("Return code was %d, expected 0", rc)
	}
	if out != "bumpversion 0.1.0\n" {
		t.Error(out)
	}
}

func TestBadFilename(t *testing.T) {
	rc, out := commandOutput([]string{"0.0.1", "**"})
	if rc != 1 {
		t.Errorf("Return code was %d, expected 1\n%v", rc, out)
	}
	rc, out = commandOutput([]string{"0.0.1", "beyzjVQLUaRrjpStWqTN.nonexistent"})
	if rc != 1 {
		t.Errorf("Return code was %d, expected 1\n%v", rc, out)
	}
}

func TestBadSemver(t *testing.T) {
	testBeforeAfter(t, "", "1.5", goFile, goFileAfter, ".go", 1)
	testBeforeAfter(t, "", "0.13.0-dev", goFile, goFileAfter, ".go", 1)
	testBeforeAfter(t, "", "latest", goFile, goFileAfter, ".go", 1)
}

func TestGo(t *testing.T) {
	testBeforeAfter(t, "", "1.3.0", goFile, goFileAfter, ".go", 0)
}

func TestMarkdown(t *testing.T) {
	testBeforeAfter(t, "", "0.15.1", md, mdAfter, ".md", 0)
}

func TestPython(t *testing.T) {
	testBeforeAfter(t, "", "0.13.2", python, pythonAfter, ".py", 0)
}

func TestReStructuredText(t *testing.T) {
	testBeforeAfter(t, "0.13.0-dev", "0.13.2", rst, rstAfter, ".rst", 0)
}

func TestSetup(t *testing.T) {
	testBeforeAfter(t, "", "0.0.3", setup, setupAfter, ".py", 0)
}

func TestUserdata(t *testing.T) {
	testBeforeAfter(t, "", "2.22.1", userdata, userdataAfter, "", 0)
}

func testBeforeAfter(t *testing.T, current, version, before, after, suffix string, code int) {
	dir, err := ioutil.TempDir("/tmp", "bumpver_test")
	if err != nil {
		t.Error(err)
	}
	filename := path.Join(dir, "testfile"+suffix)
	if err := ioutil.WriteFile(filename, []byte(before), 0644); err != nil {
		t.Error(err)
	}
	defer os.Remove(filename)
	// test bumpversion against it
	commands := []string{version, filename}
	if current != "" {
		commands = []string{"--from=" + current, version, filename}
	}
	rc, out := commandOutput(commands)
	if rc != code {
		t.Errorf("Return code was %d, expected %d", rc, code)
	}
	if code != 0 {
		return
	}
	if out != fmt.Sprintf("Bumped %s\n", filename) {
		t.Error(out)
	}
	result, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error(err)
	}
	if string(result) != after {
		t.Errorf("File contents not updated for %s:\n%s", version, string(result))
	}
}

func commandOutput(args []string) (returnCode int, output string) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rc := Command(args)

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	w.Close()
	os.Stdout = old
	output = <-outC

	return rc, output
}

var goFile = `
package version

const Version = "1.2.11"
`

var goFileAfter = `
package version

const Version = "1.3.0"
`

var python = `
#!/usr/bin/env python

from docopt import docopt

__version__ = '0.13.1'
`

var pythonAfter = `
#!/usr/bin/env python

from docopt import docopt

__version__ = '0.13.2'
`

var setup = `
#!/usr/bin/env python

"""Install the Deis command-line client."""

try:
    APACHE_LICENSE = open('LICENSE').read()
except IOError:
    APACHE_LICENSE = 'See http://www.apache.org/licenses/LICENSE-2.0'


setup(name='deis',
      version='0.0.2',
      description='Command-line Client for Deis, the open PaaS',
      classifiers=[
          'Development Status :: 4 - Beta',
          'Environment :: Console',
          'Intended Audience :: Developers',
          'Programming Language :: Python :: 2.7',
      ],
      install_requires=[
          'docopt==0.6.2', 'python-dateutil==2.2', 'requests==2.3.0', 'termcolor==1.1.0'
      ],
      **KWARGS)
`

var setupAfter = `
#!/usr/bin/env python

"""Install the Deis command-line client."""

try:
    APACHE_LICENSE = open('LICENSE').read()
except IOError:
    APACHE_LICENSE = 'See http://www.apache.org/licenses/LICENSE-2.0'


setup(name='deis',
      version='0.0.3',
      description='Command-line Client for Deis, the open PaaS',
      classifiers=[
          'Development Status :: 4 - Beta',
          'Environment :: Console',
          'Intended Audience :: Developers',
          'Programming Language :: Python :: 2.7',
      ],
      install_requires=[
          'docopt==0.6.2', 'python-dateutil==2.2', 'requests==2.3.0', 'termcolor==1.1.0'
      ],
      **KWARGS)
`

var userdata = `
#cloud-config
---
coreos:
  etcd:
    # generate a new token for each unique cluster from https://discovery.etcd.io/new
    # uncomment the following line and replace it with your discovery URL
    # discovery: https://discovery.etcd.io/12345693838asdfasfadf13939923
    addr: $private_ipv4:4001
    peer-addr: $private_ipv4:7001
    # give etcd more time if it's under heavy load - prevent leader election thrashing
    peer-election-timeout: 2000
    # heartbeat interval should ideally be 1/4 or 1/5 of peer election timeout
    peer-heartbeat-interval: 500
write_files:
  - path: /etc/deis-release
    content: |
      DEIS_RELEASE=2.22.0
  - path: /etc/profile.d/nse-function.sh
    permissions: '0755'
    content: |
      function nse() {
        sudo nsenter --pid --uts --mount --ipc --net --target $(docker inspect --format="{{ .State.Pid }}" $1)
      }
`

var userdataAfter = `
#cloud-config
---
coreos:
  etcd:
    # generate a new token for each unique cluster from https://discovery.etcd.io/new
    # uncomment the following line and replace it with your discovery URL
    # discovery: https://discovery.etcd.io/12345693838asdfasfadf13939923
    addr: $private_ipv4:4001
    peer-addr: $private_ipv4:7001
    # give etcd more time if it's under heavy load - prevent leader election thrashing
    peer-election-timeout: 2000
    # heartbeat interval should ideally be 1/4 or 1/5 of peer election timeout
    peer-heartbeat-interval: 500
write_files:
  - path: /etc/deis-release
    content: |
      DEIS_RELEASE=2.22.1
  - path: /etc/profile.d/nse-function.sh
    permissions: '0755'
    content: |
      function nse() {
        sudo nsenter --pid --uts --mount --ipc --net --target $(docker inspect --format="{{ .State.Pid }}" $1)
      }
`

var rst = `
.. code-block:: console

    $ pip install docopt==0.6.2 python-dateutil==2.2 requests==2.3.0 termcolor==1.1.0
    $ sudo ln -fs $(pwd)/client/deis.py /usr/local/bin/deis
    $ deis
    Usage: deis <command> [<args>...]

If you don't have Python_ installed, you can download a binary executable
version of the Deis client for Mac OS X, Linux amd64, or Windows:

    - https://s3-us-west-2.amazonaws.com/opdemand/deis-0.13.0-dev-darwin.tgz
    - https://s3-us-west-2.amazonaws.com/opdemand/deis-0.13.0-dev-linux.tgz
`

var rstAfter = `
.. code-block:: console

    $ pip install docopt==0.6.2 python-dateutil==2.2 requests==2.3.0 termcolor==1.1.0
    $ sudo ln -fs $(pwd)/client/deis.py /usr/local/bin/deis
    $ deis
    Usage: deis <command> [<args>...]

If you don't have Python_ installed, you can download a binary executable
version of the Deis client for Mac OS X, Linux amd64, or Windows:

    - https://s3-us-west-2.amazonaws.com/opdemand/deis-0.13.2-darwin.tgz
    - https://s3-us-west-2.amazonaws.com/opdemand/deis-0.13.2-linux.tgz
`

var md = `
# Deis

Deis (pronounced DAY-iss) is an open source PaaS that makes it easy to deploy and manage applications on your own servers. Deis builds upon [Docker](http://docker.io/) and [CoreOS](http://coreos.com) to provide a lightweight PaaS with a [Heroku-inspired](http://heroku.com) workflow.

[![Current Release](http://img.shields.io/badge/release-v0.12.0-blue.svg)](https://github.com/deis/deis/releases/tag/v0.12.0)

![Deis Graphic](https://s3-us-west-2.amazonaws.com/deis-images/deis-graphic.png)
`

var mdAfter = `
# Deis

Deis (pronounced DAY-iss) is an open source PaaS that makes it easy to deploy and manage applications on your own servers. Deis builds upon [Docker](http://docker.io/) and [CoreOS](http://coreos.com) to provide a lightweight PaaS with a [Heroku-inspired](http://heroku.com) workflow.

[![Current Release](http://img.shields.io/badge/release-v0.15.1-blue.svg)](https://github.com/deis/deis/releases/tag/v0.15.1)

![Deis Graphic](https://s3-us-west-2.amazonaws.com/deis-images/deis-graphic.png)
`
