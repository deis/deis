# Contributing

## Issues

When opening issues, try to be as clear as possible when describing the bug or feature request.
Tag the issue accordingly.

## Pull Requests

To hack on targz:

1. Install as usual (`go get github.com/walle/targz`)
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Ensure everything works and the tests pass (see below)
4. Commit your changes (`git commit -am 'Add some feature'`)

Contribute upstream:

1. Fork targz on GitHub
2. Add your remote (`git remote add fork git@github.com:myuser/repo.git`)
3. Push to the branch (`git push fork my-new-feature`)
4. Create a new Pull Request on GitHub

For other team members:

1. Install as usual (`go get github.com/walle/targz`)
2. Add your remote (`git remote add fork git@github.com:myuser/repo.git`)
3. Pull your revisions (`git fetch fork; git checkout -b my-new-feature fork/my-new-feature`)

Notice: Always use the original import path by installing with `go get`.

## Testing

To run the test suite use the command

    $ go test -cover

The tests will write a file structure to a temporary directory on your disk.