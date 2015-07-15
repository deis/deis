:title: Test Plan
:description: Details the strategy used to verify the Deis platform

.. _test_plan:

Test Plan
=========

Identifier
----------

This document's identifier is **OPD-MTP-1.0.1**. The identifier changes with
any significant revision to the document.


References
----------

- "Testing Deis": http://docs.deis.io/en/v1.0.1/contributing/testing/
- "Changes Must Include Tests": http://docs.deis.io/en/v1.0.1/contributing/standards/#include-tests
- "Release Schedule": http://docs.deis.io/en/v1.0.1/contributing/schedule/


Introduction
------------

This document describes the master test plan for the Deis open source PaaS
software platform. It is a *living document*, continually updated to reflect the
current practice and goals of software quality assurance for Deis.

Deis has unit_, functional_, integration_, and acceptance_ tests that cover
essential product functionality. The Deis project relies on contributors to be
responsible by exercising these tests themselves, and by including tests when
proposing changes. Deis is tested and validated by a 24/7
`continuous integration`_ platform, supplemented by intelligent manual testing.


Test items
----------

Within the scope of this master test plan are these items:

- The Deis project codebase at https://github.com/deis/deis
- The Docker containers that constitute Deis
- The assembled Deis platform on a CoreOS cluster
- The HTML documentation set for Deis
- Binary installers for the ``deis`` CLI hosted at AWS S3
- Binary installers for the ``deisctl`` CLI hosted at AWS S3
- Hosted HTML documentation updates at http://docs.deis.io/
- Docker images hosted at https://registry.hub.docker.com/repos/deis/


.. _features_to_be_tested:

Features to be Tested
---------------------

At a high level, the overall features of the platform that are tested are:

- The Deis platform can be installed on a new CoreOS vagrant cluster
- The Deis platform can be upgraded from a recent release to the current one
- Users can register with Deis and create and deploy applications
- Deis can build and scale a variety of Heroku-style and Dockerfile-based apps
- Users can grant and revoke application access to other users


Features not to be Tested
-------------------------

While these features are effectively covered in ad-hoc testing and by existing
customer usage, they are not specifically tested as part of the test plan yet.

- The Deis platform can be installed on an existing CoreOS cluster
- The Deis platform survives the loss of one or more nodes
- The Deis platform can run for several weeks without interruption
- The Deis platform handles stressful multi-user workloads

These features are not included in the test plan currently due to resource
limitations. Future test automation will move these features into the
"to be tested" section.


Approach
--------

Deis' test plan relies on extensive test automation, supplemented by spot
testing by responsible developers. Continuous integration tests ensure that
the platform functions and regressions are not introduced, and focused manual
testing is also relied upon for acceptance testing a product release.

Developers are expected to have run the same tests locally that will be run
for them by continuous integration, specifically the test-integration.sh_
script. This will execute documentation tests, unit tests, functional tests,
and then an overall acceptance test against a Vagrant cluster.

As changes are incorporated into Deis and the team plans a product release,
maintainers begin acceptance tests against other cloud providers, following the
instructions exactly as provided to users. When these platform-specific tests
have passed, a final validation test occurs in continuous integration against
a tagged codebase. If it succeeds, the release has passed.

Nightly jobs run a subset of test-integration tests against the released CLI
installers and the current codebase to detect regressions in product behavior.


Item Pass/Fail Criteria
-----------------------

Integration testing has passed when no failures occurred in the
test-integration.sh_ script as run by continuous integration. (Each proposed
change also requires review and approval by two project maintainers before it
is actually merged into the codebase, see :ref:`merge_approval`.)

Acceptance testing has passed when no failures occurred in the "test-latest.sh"
script as run by continuous integration in the "test-latest" job, and when
maintainers have completed spot testing successfully, as defined by local
release criteria.


Suspension Criteria and Resumption Requirements
-----------------------------------------------

Suspension of automated testing occurs when any failure arises in any part of
the test suite. There are no optional or weighted failures: everything
must pass. Suspension of manual testing occurs when a failure arises that
makes further testing unpredictable or of limited value.

Resumption of testing occurs when a test failure has been addressed and fixed
such that it is reasonable to assume tests may pass again.


Test Deliverables
-----------------

- Detailed testing logs as generated by CI jobs.
  See https://ci.deis.io/job/test-master/85/console for an example.


Remaining Test Tasks
--------------------

For acceptance testing, test deliverables are generated by a maintainer starting
the CLI and test-master jobs when appropriate.


Environmental needs
-------------------

Testing requires a Linux or Mac OS X host capable of running VirtualBox and
Vagrant with good network connectivity. Specific environmental needs are
outlined in the setup-node.sh_ script, which should be kept up-to-date with
current needs.


Staffing and training needs
---------------------------

N/A


Responsibilities
----------------

A maintainer designated as "QA Lead" for an acceptance test process has the
responsibility to execute the test task of starting appropriate CI jobs. The
QA Lead is also tasked with overseeing manual testing activities executed by
others.

For a patch or minor release, the QA Lead may decide not to execute all aspects
of acceptance testing.

The QA Lead may also execute clerical tasks associated with a release as
described in the :ref:`releases` documentation.


Schedule
--------

As we describe an ongoing, evolving test plan here, there is no fixed project
schedule to address, just a repeatable process.

Deis releases early and often. The consequences of a failure in the test process
described here are a delay to an expected release date and the restart of the
test process once the failure has been addressed.


Risks and Contingencies
-----------------------

Automated tests do not yet extend to all cloud providers, and it is possible
that manual testing could miss something. We will address this by adding AWS_
and other testing flavors soon.

Resources are limited, and contention between development needs and testing
needs has the potential to slow down the quality assurance process.


Approvals
---------

The Deis maintainer team as a whole approves this document through our normal
pull request and merge approval process. Comments and additions will be made as
pull requests against this documentation.


.. _unit: http://en.wikipedia.org/wiki/Unit_testing
.. _functional: http://en.wikipedia.org/wiki/Functional_testing
.. _integration: http://en.wikipedia.org/wiki/Integration_testing
.. _acceptance: http://en.wikipedia.org/wiki/Acceptance_testing
.. _`continuous integration`: http://en.wikipedia.org/wiki/Continuous_integration
.. _`source code`: https://github.com/deis/deis
.. _test-integration.sh: https://github.com/deis/deis/blob/master/tests/bin/test-integration.sh
.. _setup-node.sh: https://github.com/deis/deis/blob/master/tess/bin/setup-node.sh
.. _AWS: http://aws.amazon.com/
