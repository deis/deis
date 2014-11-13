# bumpver

`bumpver` is a go command-line tool that changes [semantic version](http://semver.org/)
strings in source and doc files. It was created to help with [Deis](http://deis.io/)
releases.

## Build

```console
$ make test build
$ ./bumpver --help
Updates the current semantic version number to a new one in the specified
source and doc files.

Usage:
  bumpver [-f <current>] <version> <files>...

Options:
  -f --from=<current>  An explicit version string to replace. Otherwise, use
                       the first semantic version found in the first file.
```

## Usage

Let's use `bumpver` to update the deis codebase for a new release, 0.13.3:

```console
$ make -C contrib/bumpver/ test build
$ # update from an explicit (bad) version string
$ ./contrib/bumpver/bumpver -f latest 0.13.3 contrib/coreos/user-data
$ ./contrib/bumpver/bumpver -f 0.13.0-dev 0.13.3 \
    version/version.go \
    client/deis.py \
    client/setup.py \
    deisctl/deis-version \
    deisctl/deisctl.go \
    deisctl/README.md \
    controller/deis/__init__.py \
    README.md
$ # update from the first semver string found
$ # this type of command should now be enough to bump everything
$ ./contrib/bumpver/bumpver 0.14.0 \
    version/version.go \
    client/deis.py \
    client/setup.py \
    deisctl/deis-version \
    deisctl/deisctl.go \
    deisctl/README.md \
    contrib/coreos/user-data.example \
    controller/deis/__init__.py \
    README.md
```

Of course, you should **always** check the changes with `git diff` before committing
anything to version control. You can also check the
[Release Checklist](http://docs.deis.io/en/latest/contributing/releases/) for the
most up-to-date list of files to bump.

Please add any issues you find with this software to the
[Deis project](https://github.com/deis/deis/issues).
