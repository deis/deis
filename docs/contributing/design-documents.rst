:title: Design Documents
:description: Information necessary for a feature design document.

.. _design-documents:

Design Documents
================

Before submitting a pull request which will change the behavior of Deis significantly,
such as a new feature or major refactoring, contributors should first open
an issue representing a design document.

Goals
-----

Design documents help ensure project contributors:

* Involve stakeholders as early as possible in a feature's development
* Ensure code changes accomplish the original motivations and design goals
* Establish clear acceptance criteria for a feature or change
* Enforce test-driven design methodology and automated test coverage

Contents
--------

Design document issues should be named ``Design Doc: <change description>`` and
contain the following sections.

Goal
^^^^

This section should briefly describe the proposed change and the motivations
behind it. Tests will be written to ensure this design goal is met by
the change.

This section should also reference a separate GitHub issue tracking
the feature or change, which will typically be assigned to a release milestone.

Code Changes
^^^^^^^^^^^^

This section should detail the code changes necessary to accomplish the change,
as well as the proposed implementation. This should be as detailed as necessary to
help reviewers understand the change.

Tests
^^^^^

All changes should be covered by automated tests, either unit or integration tests
(ideally both). This section should detail how tests will be written to validate
that the change accomplishes the design goals and doesn't introduce any regressions.

If a change cannot be sufficiently covered by automated testing, the design
should be reconsidered. If there is no test coverage whatsoever for an affected
section of code, a separate issue should be filed to integrate automated testing
with that section of the codebase.

The tests described here also form the acceptance criteria for the change, so
that when it's completed maintainers can merge the pull request after confirming
the tests pass CI.

Approval
--------

A design document follows the same :ref:`merge_approval` review process as final
PRs do, and maintainers will take extra care to ensure that any stakeholders for
the change are included in the discussion and review of the design document.

Once the design is accepted, the author can complete the change and submit a pull
request for review. The pull request should close both the design document for
the change as well as any issues that either track the issue or are closed as a
result of the change.

See :ref:`standards` for more information on pull request and commit message formatting.
