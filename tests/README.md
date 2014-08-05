# Deis Tests

This directory contains a [Go](http://golang.org/) package with integration
tests for the [Deis](http://deis.io/) open source PaaS.

[![GoDoc](https://godoc.org/github.com/deis/deis/tests?status.svg)](https://godoc.org/github.com/deis/deis/tests)

**NOTE**: These integration tests are targeted for use in Deis'
[continuous integration system](http://ci.deis.io/). The tests currently assume
they are targeting a freshly provisioned Deis cluster. **Don't** run the
integration tests on a Deis installation with existing users; the tests will
fail and could overwrite data.

## Test Setup

Check out [Deis' source code](https://github.com/deis/deis) into the `$GOPATH`:

```console
$ go get -u -v github.com/deis/deis
$ cd $GOPATH/src/github.com/deis/deis/tests
```

Provision a Deis cluster as usual, and ensure that a matching `deis`
command-line client is available in your `$PATH`.

Create two SSH keys:

```console
$ ssh-keygen -q -t rsa -f ~/.ssh/deis -N '' -C deis
$ ssh-keygen -q -t rsa -f ~/.ssh/deiskey -N '' -C deiskey
```

The first key `deis` is used for authentication against Deis by the `test` user
who runs the integration tests. The second `deiskey` is used only for testing
`deis keys:add`, `deis keys:list`, and related commands.

## Test Execution

Run all the integration tests:

```console
$ make test-full
```

Or run just the [smoke test](http://www.catb.org/jargon/html/S/smoke-test.html):

```console
$ make test-smoke
```

## Customizing Test Runs

These environment variables can be set to customize the test run:
* `DEIS_TEST_AUTH_KEY` - SSH key used to register with the Deis controller
  (default: `~/.ssh/deis`)
* `DEIS_TEST_SSH_KEY` - SSH key used to login to the controller machine
  (default: `~/.vagrant.d/insecure_private_key`)
* `DEIS_TEST_DOMAIN` - the domain to use for testing
  (default: `local.deisapp.com`)
* `DEIS_TEST_HOSTS` - comma-separated list of IPs for nodes in the cluster,
  should be internal IPs for cloud providers (default: `172.17.8.100`)
* `DEIS_TEST_APP` - name of the
  [Deis example app](https://github.com/deis?query=example-) to use, which is
  cloned from GitHub (default: random)
