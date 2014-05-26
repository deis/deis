# Deis Integration Tests

This directory contains integration tests to be run against an existing Deis cluster.

To run all tests (from the project root):

```console
$ vagrant up
$ make pull run
$ make test
```

Tests can also be run manually:

```console
$ go test -v ./test/integration/deis_test.go
```

The test environment uses several environment variables, which can be set to customize the run:

DEIS_TEST_KEY:
    SSH key used to login to the controller machine

DEIS_TEST_HOSTNAME:
    hostname which resolves to the controller host

DEIS_TEST_HOSTS:
    comma-separated list of IPs for nodes in the cluster.
    These should be internal IPs for cloud providers

DEIS_TEST_APP_URL:
    URL of the Deis example app to use, which is cloned from GitHub
