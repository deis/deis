# Deis integration testing

This directory contains a Go package which will run integration tests on an
existing Deis cluster.

Repo should be properly checked out into your GOPATH
go get github.com/deis/deis

To run all tests:

```console
$ go test -v -tags integration -timeout 50m ./...
```

The test environment uses several environment variables, which can be set to customize the run:
* `DEIS_TEST_AUTH_KEY` - SSH key used to register with the Deis controller -- will be generated if it doesn't exist (default: `~/.ssh/deis`)
* `DEIS_TEST_KEY` - SSH key used to login to the controller machine (default: `~/.vagrant.d/insecure_private_key`)
* `DEIS_TEST_HOSTNAME` - hostname which resolves to the controller host (default: `local.deisapp.com`)
* `DEIS_TEST_HOSTS` - comma-separated list of IPs for nodes in the cluster -- should be internal IPs for cloud providers (default: `172.17.8.100`)
* `DEIS_TEST_APP` - name of the Deis example app to use, which is cloned from GitHub (default: `example-ruby-sinatra`)
