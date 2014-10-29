:title: Coding Standards
:description: How to propose changes to the Deis codebase.


.. _standards:

PR Checklist
============

Proposed changes to Deis are made as GitHub `pull requests`_.

Please make sure your PR follows this checklist:

1. `Single Issue`_
2. `Include Tests`_
3. `Include Docs`_
4. `Code Standards`_
5. `Commit Style`_


Single Issue
------------

When fixing or implementing a GitHub issue, resist the temptation to refactor
nearby code or to fix that potential bug you noticed. Instead, open a new
`pull request`_ just for that change.

It's hard to reach agreement on the merit of a PR when it isn't focused. Keeping
concerns separated allows pull requests to be tested, reviewed, and merged
more quickly.

Squash and rebase the commit or commits in your pull request into logical units
of work with ``git``. Include tests and documentation changes in the same commit,
so that a revert would remove all traces of the feature or fix.

Most pull requests will reference a `GitHub issue`_. In the PR description--not
in the commit itself--include a line such as "Closes #1234." The issue referenced
will then be closed when your PR is merged.


Include Tests
-------------

While working on local code changes, run Deis' tests:

.. code-block:: console

    $ export DEV_REGISTRY=192.168.59.103:5000 HOST_IPADDR=192.168.59.103
    $ ./tests/bin/test-integration.sh

    >>> Preparing test environment <<<

    DEIS_ROOT=/Users/matt/Projects/src/github.com/deis/deis
    DEIS_TEST_APP=example-go
    ...

You can run subsets of the tests:

.. code-block:: console

    $ make -C docs/ test
    $ make -C controller/ test-unit

Be sure your proposed changes pass all of ``./tests/bin/test-integration``
on your workstation before submitting a PR.

See ``tests/README.md`` in the code for more information.


Include Docs
------------

Any change to Deis that could affect a user's experience also needs a change or
addition to the relevant documentation. Deis generates the HTML documentation
hosted at http://docs.deis.io/ from the text markup sources in the
``docs/`` directory.

See ``docs/README.md`` in the code for more information.


Code Standards
--------------

Deis is a Go_ and Python_ project. For both languages, we agree with
`The Zen of Python`_, which emphasizes simple over clever. Readability counts.

Go code should always be run through ``gofmt`` on the default settings. Lines
of code may be up to 99 characters long. Documentation strings and tests are
required for all public methods. Use of third-party go packages should be
minimal, but when doing so, vendor code into the Deis package with the
godep_ tool.

Python code should always adhere to PEP8_, the python code style guide, with
the exception that lines of code may be up to 99 characters long. Docstrings and
tests are required for all public methods, although the flake8_ tool used by
Deis does not enforce this.


.. _commit_style_guide:

Commit Style
------------

``git commit`` messages must follow this format::

    {type}({scope}): {subject}
    <BLANK LINE>
    {body}
    <BLANK LINE>
    {footer}

Example
"""""""

::

    feat(logger): add frobnitz pipeline spout discovery

    Introduces a FPSD component compatible with the industry standard for
    spout discovery.

    BREAKING CHANGE: Fixing the buffer overflow in the master subroutine
        required losing compatibility with the UVEX-9. Any UVEX-9 or
        umVEX-8 series artifacts will need to be updated to umVX format
        with the consortium or vendor toolset.


Subject Line
""""""""""""

The first line of a commit message is its subject. It contains a brief
description of the change, no longer than 50 characters.

These {types} are allowed:

- **feat** -> feature
- **fix** -> bug fix
- **docs** -> documentation
- **style** -> formatting
- **ref** -> refactoring code
- **test** -> adding missing tests
- **chore** -> maintenance

The {scope} specifies the location of the change, such as "controller,"
"Dockerfiles," or ".gitignore". The {subject} should use an imperative,
present-tense verb: "change," not "changes" or "changed." Don't
capitalize the verb or add a period (.) at the end of the subject line.

Message Body
""""""""""""

Separate the message body from the subject with a blank line. The body
can have lines up to 72 characters long. It includes the motivation for the
change and points out differences from previous behavior. The body and
the footer should be written as full sentences.

Message Footer
""""""""""""""

Separate a footer from the message body with a blank line. Mention any
breaking change along with the justification and migration notes. If the
changes cannot be tested by Deis' test scripts, include specific instructions
for manual testing.


Merge Approval
--------------

Deis maintainers add "**LGTM**" (Looks Good To Me) or an equivalent comment
to indicate that a PR is acceptable. Any code change--other than
a simple typo fix or one-line documentation change--requires at least two
maintainers to accept it.

If the PR is from a Deis maintainer, then he or she should be the one to close
it. This keeps the commit stream clean and gives the maintainer the benefit of
revisiting the PR before deciding whether or not to merge the changes.


.. _Python: http://www.python.org/
.. _Go: http://golang.org/
.. _godep: https://github.com/tools/godep
.. _flake8: https://pypi.python.org/pypi/flake8/
.. _PEP8: http://www.python.org/dev/peps/pep-0008/
.. _`The Zen of Python`: http://www.python.org/dev/peps/pep-0020/
.. _`pull request`: https://github.com/deis/deis/pulls
.. _`pull requests`: https://github.com/deis/deis/pulls
.. _`GitHub issue`: https://github.com/deis/deis/issues
