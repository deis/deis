:title: Coding Standards
:description: Deis project coding standards. Contributors to Deis should feel welcome to make changes to any part of the codebase.

.. _standards:

Coding Standards
================

Deis is a Python_ and Go_ project.

We chose Go_ because it is simple, reliable, and efficient. These are
values Deis shares. Go also excels at programming systems-level tasks,
with powerful and simple networking, concurrency, and testing facilities
included. Our coding standards and recommendations for Go code in the
Deis project are evolving, but will be added to this document soon.

We chose Python_ over other compelling languages because it is
widespread, well-documented, and friendly to a large number of
developers. Source code benefits from many eyes upon it.

`The Zen of Python`_ emphasizes simple over clever, and we agree.
Readability counts. Deis also aims for complete test coverage.

Contributors to Deis should feel welcome to make changes to any part
of the codebase. To create a proper GitHub pull request for inclusion
into the official repository, your code must pass two tests:

- :ref:`make_flake8`
- :ref:`make_coverage`


.. _make_flake8:

``make flake8``
---------------

`flake8`_ is a helpful command-line tool that combines the output of
`pep8 <pep8_tool_>`_, `pyflakes`_, and `mccabe`_.

.. code-block:: console

    $ make -C controller flake8
    flake8

No output, as above, means ``flake8`` found no errors. If errors
are reported, fix them in your source code and try ``flake8`` again.

The Deis project adheres to `PEP8`_, the python code style guide,
with the exception that we allow lines up to 99 characters in length.
Docstrings and tests are also required for all public methods, although
``flake8`` does not enforce this.

Default settings for ``flake8`` are in the ``[flake8]`` section of the
setup.cfg file in the project root.


.. _make_coverage:

``make coverage``
-----------------

Once your code passes the style checker, run the test suite and
ensure that everything passes and that code coverage has not declined.

.. code-block:: console

    $ make -C controller coverage
    coverage run manage.py test --noinput api web
    WARNING Cannot synchronize with etcd cluster
    Creating test database for alias 'default'...
    ...............................................
    ----------------------------------------------------------------------
    Ran 47 tests in 47.768s

    OK
    Destroying test database for alias 'default'...
    coverage html

If a test fails, fixing it is obviously the first priority. And if you
have introduced new code, it must be accompanied by unit tests.

In the example above, all tests passed and ``coverage`` created a report
of what code was exercised while the tests were running. Open the file
``htmlcov/index.html`` under the project's root and ensure that the
overall coverage percentage has not receded as a result of your
changes. Current test coverage can be found here:

.. image:: https://coveralls.io/repos/deis/deis/badge.png?branch=master
    :target: https://coveralls.io/r/deis/deis?branch=master
    :alt: Coverage Status


.. _pull_request:

Pull Request
------------

Now create a GitHub `pull request`_ with a description of what your code
fixes or improves.

Before the pull request is merged, make sure that you squash your
commits into logical units of work using ``git rebase -i`` and
``git push -f``. Include documentation changes in the same commit,
so that a revert would remove all traces of the feature or fix.

Commits that fix or close an issue should include a reference like
*Closes #XXX* or *Fixes #XXX* in the commit message. Doing so will
automatically close the `GitHub issue`_ when the pull request is merged.

Merge Approval
--------------

Deis maintainers add "**LGTM**" (Looks Good To Me) in code
review comments to indicate that a PR is acceptable. Any code change--other than
a simple typo fix or one-line documentation change--requires at least two of
Deis' maintainers to accept the change in this manner before it can be merged.
If the PR is from a Deis maintainer, then he or she should be the one to merge
it. This is for cleanliness in the commit stream as well as giving the
maintainer the benefit of adding more fixes or commits to a PR before the
merge.

.. _Python: http://www.python.org/
.. _Go: http://golang.org/
.. _flake8: https://pypi.python.org/pypi/flake8/
.. _pep8_tool: https://pypi.python.org/pypi/pep8/
.. _pyflakes: https://pypi.python.org/pypi/pyflakes/
.. _mccabe: https://pypi.python.org/pypi/mccabe/
.. _PEP8: http://www.python.org/dev/peps/pep-0008/
.. _`The Zen of Python`: http://www.python.org/dev/peps/pep-0020/
.. _`pull request`: https://github.com/deis/deis/pulls
.. _`GitHub issue`: https://github.com/deis/deis/issues


.. _commit_style_guide:

Commit Style Guide
------------------

There are several reasons why we try to follow a specific style guide for commits:

- it allows us to recognize unimportant commits like formatting
- it provides better information when browsing the git history

Recognizing Unimportant Commits
```````````````````````````````

These commits are usually just formatting changes like adding/removing spaces/empty lines,
fixing indentation, or adding comments. So when you are looking for some change in the
logic, you can ignore these commits - there's no logic change inside this commit.

When bisecting, you can ignore these by running:

.. code-block:: console

    git bisect skip $(git rev-list --grep irrelevant <good place> HEAD)

Providing more Information when Browsing the History
````````````````````````````````````````````````````

This adds extra context to our commit logs. Look at these messages (taken from the last
few AngularJS commits):

- Fix small typo in docs widget (tutorial instructions)
- Fix test for scenario.Application - should remove old iframe
- docs - various doc fixes
- docs - stripping extra new lines
- Replaced double line break with single when text is fetched from Google
- Added support for properties in documentation

All of these messages try to specify where the change occurs, but they don’t share any
convention. Now look at these messages:

- fix comment stripping
- fixing broken links
- Bit of refactoring
- Check whether links do exist and throw exception
- Fix sitemap include (to work on case sensitive linux)

Are you able to guess what’s inside each commit diff?

It's true that you can find this information by checking which files had been changed, but
that’s slow. When looking in the git history, we can see that all of the developers are
trying to specify where the change takes place, but the message is missing a convention.
Cue commit message formatting entrance stage left.

Format of the Commit Message
````````````````````````````

.. code-block:: console

    {type}({scope}): {subject}
    <BLANK LINE>
    {body}
    <BLANK LINE>
    {footer}

Any line of the commit message cannot be longer than 72 characters, with the subject
line limited to 50 characters. This allows the message to be easier to read on github
as well as in various git tools. Lines of source code prefixed with four spaces are exempt
from this rule.

Subject Line
""""""""""""

The subject line contains a succinct description of the change to the logic.

The allowed {types} are as follows:

- feat -> feature
- fix -> bug fix
- docs -> documentation
- style -> formatting
- ref -> refactoring code
- test -> adding missing tests
- chore -> maintenance

The {scope} can be anything specifying place of the commit change e.g. the controller,
the client, the logger, etc.

The {subject} needs to use imperative, present tense: “change”, not “changed” nor
“changes”. The first letter should not be capitalized, and there is no dot (.) at the end.

Message Body
""""""""""""

Just like the {subject}, the message {body} needs to be in the present tense, and includes
the motivation for the change, as well as a contrast with the previous behavior.

Message Footer
""""""""""""""

All breaking changes need to be mentioned in the footer with the description of the
change, the justification behind the change and any migration notes required. Any methods
that maintainers can use to test these changes should be placed in the footer as well. For
example:

.. code-block:: console

    TESTING: to test this change, bring up a new cluster and run the following
    when the controller comes online:

        $ vagrant ssh -c "curl localhost:8000"

    you should see an HTTP response from the controller.

    BREAKING CHANGE: the controller no longer listens on port 80. It now listens on
    port 8000, with the router redirecting requests on port 80 to the controller. To
    migrate to this change, SSH into your controller and run:

        $ docker kill deis-controller
        $ docker rm deis-controller

    and then restart the controller on port 8000:

        $ docker run -d -p 8000:8000 -e ETCD=<etcd_endpoint> -e HOST=<host_ip> \
        -e PORT=8000 -name deis-controller deis/controller

    now you can start the proxy component by running:

        $ docker run -d -p 80:80 -e ETCD=<etcd_endpoint> -e HOST=<host_ip> -e PORT=80 \
        -name deis-router deis/router

    the router should then start proxying requests from port 80 to the controller.

Referencing Issues
""""""""""""""""""

Closed bugs should be listed on a separate line in the footer prefixed with the "closes"
keyword like this:

.. code-block:: console

    closes #123

Or in the case of multiple issues:

.. code-block:: console

    closes #123, #456, #789

Examples
````````

.. code-block:: console

    feat(controller): add router component

    This introduces a new router component to Deis, which proxies requests to Deis
    components.

    closes #123

    BREAKING CHANGE: the controller no longer listens on port 80. It now listens on
        port 8000, with the router redirecting requests on port 80 to the controller. To
        migrate to this change, SSH into your controller and run:

        $ docker kill deis-controller
        $ docker rm deis-controller

        and then restart the controller on port 8000:

        $ docker run -d -p 8000:8000 -e ETCD=<etcd_endpoint> -e HOST=<host_ip> \
        -e PORT=8000 -name deis-controller deis/controller

        now you can start the proxy component by running:

        $ docker run -d -p 80:80 -e ETCD=<etcd_endpoint> -e HOST=<host_ip> -e PORT=80 \
        -name deis-router deis/router

        The router should then start proxying requests from port 80 to the controller.
    ----------------------------------------------------------------------------------
    test(client): add unit tests for app domains

    Nginx does not allow domain names larger than 128 characters, so we need to make
    sure that we do not allow the client to add domains larger than 128 characters.
    A DomainException is raised when the domain name is larger than the maximum
    character size.

    closes #392

Signing your Commits
""""""""""""""""""""

Unlike docker_ or other parallel projects, we do not require that you need to
sign your commits. Submitting a PR certifies that you have read and agree to the terms in
the DCO.

.. _docker: https://docker.com/

Forcing no Build for TravisCI
"""""""""""""""""""""""""""""

If you're committing a PR that is just a small typo fix or a README change, you can force
TravisCI to not build your commit by adding [skip ci] below the message body. For example:

.. code-block:: console

    fix(README): typo

    It's spelled tomato, not tomatoe.

    [skip ci]
