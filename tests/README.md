# Deis Tests

This directory contains a [Go](http://golang.org/) package with integration
tests for the [Deis](http://deis.io/) open source PaaS.

Please refer to [Testing Deis](http://docs.deis.io/en/latest/contributing/testing/)
for help with running these tests.

[![Testing Deis](https://readthedocs.org/projects/deis/badge/)](http://docs.deis.io/en/latest/contributing/testing/)

**NOTE**: These integration tests are targeted for use in Deis'
[continuous integration system](https://ci.deis.io/). The tests currently assume
they are targeting a freshly provisioned Deis cluster. **Don't** run the
integration tests on a Deis installation with existing users; the tests will
fail and could overwrite data.
