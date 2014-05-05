# Deis integration testing
This directory includes a Rakefile which will run integration tests on an existing Deis cluster.

To run all tests:

```console
$ bundle install
$ bundle exec rake
```

The namespaces `setup`, `tests`, and `cleanup` are defined. The default task runs `setup:all`, `tests:all`, `cleanup:all` and then exits.

Namespaces can also be run manually:

```console
$ bundle exec rake setup:all
```

...and so can tests:

```console
$ bundle exec rake tests:create_cluster
```

The test environment uses several environment variables, which can be set to customize the run:
* `DEIS_TEST_AUTH_KEY` - SSH key used to register with the Deis controller -- will be generated if it doesn't exist (default: `~/.ssh/deis`)
* `DEIS_TEST_KEY` - SSH key used to login to the controller machine (default: `~/.vagrant.d/insecure_private_key`)
* `DEIS_TEST_HOSTNAME` - hostname which resolves to the controller host (default: `local.deisapp.com`)
* `DEIS_TEST_HOSTS` - comma-separated list of IPs for nodes in the cluster -- should be internal IPs for cloud providers (default: `172.17.8.100`)
* `DEIS_TEST_APP` - name of the Deis example app to use, which is cloned from GitHub (default: `example-ruby-sinatra`)
