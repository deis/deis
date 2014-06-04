# How to Contribute

The Deis project is Apache 2.0 licensed and accepts contributions via Github pull
requests. This document outlines some of the conventions on commit message formatting,
contact points for developers and other resources to make getting your contribution
accepted.

# Certificate of Origin

By contributing to this project you agree to the
[Developer Certificate of Origin (DCO)][dco]. This document was created by the Linux
Kernel community and is a simple statement that you, as a contributor, have the legal
right to make the contribution.

# Support Channels

- IRC: #[deis](irc://irc.freenode.org:6667/#deis) IRC channel on freenode.org

## Getting Started

- Fork the repository on GitHub
- Read the README.md for build instructions

## Contribution Flow

This is a rough outline of what a contributor's workflow looks like:

- Create a topic branch from where you want to base your work. This is usually master.
- Make commits of logical units.
- Make sure your commit messages are in the proper format, see below
- Push your changes to a topic branch in your fork of the repository.
- Submit a pull request

Thanks for your contributions!

### Format of the Commit Message

We follow a rough convention for commit messages borrowed from CoreOS, who borrowed theirs
from AngularJS. This is an example of a commit:

```
    feat(scripts/test-cluster): add a cluster test command

    this uses tmux to setup a test cluster that you can easily kill and
    start for debugging.
```

To make it more formal it looks something like this:

```
{type}({scope}): {subject}
<BLANK LINE>
{body}
<BLANK LINE>
{footer}
```

Any line of the commit message cannot be longer than 90 characters, with the subject line
limited to 70 characters. This allows the message to be easier to read on github as well
as in various git tools.

The allowed {types} are as follows:

```
feat -> feature
fix -> bug fix
docs -> documentation
style -> formatting
refactor
test -> adding missing tests
chore -> maintenance
```

### More Details on Commits

For more details see the [commit style guide][style-guide].

[dco]: DCO
[style-guide]: http://docs.deis.io/en/latest/contributing/standards/#commit-style-guide
