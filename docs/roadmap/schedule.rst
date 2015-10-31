:title: Release Schedule
:description: When does Deis have major, minor, and patch releases?

.. _release_schedule:

Release Schedule
================

Some of the greatest assets of the Deis project are velocity and agility.
Deis changed rapidly and with relative ease during initial development.

Deis now harnesses those strengths into powering a regular, public release
cadence. From v1.0.0 onward, the Deis project will release a minor version each
month, with patch versions as needed. The project will use GitHub milestones to
communicate the content and timing of major and minor releases.

Deis releases are not feature-based, in that dates are not linked to specific
features. If a feature is merged before the release date, it is included in the
next minor or major release.

The master ``git`` branch of Deis always works. Only changes considered ready to
be released publicly are merged, and releases are made from master.

Semantic Versioning
-------------------

Deis releases comply with `semantic versioning`_, with the "public API" broadly
defined as:

- the REST API for *deis-controller*
- etcd keys and values that are publicly documented
- ``deis`` and ``deisctl`` commands and options
- essential ``Makefile`` targets
- provider scripts under ``contrib/``

Users of Deis can be confident that upgrading to a patch or to a minor release
will not change the behavior of these items in a backward-incompatible way.

Release Criteria
----------------

For any Deis release to be made publicly available, it must meet at least
these criteria:

- Passes all tests on the supported, load-balancing cloud providers
- Has no new regressions in behavior that are not considered trivial

Patch Releases
--------------

A patch release of Deis includes backwards-compatible bug fixes. Upgrading to
this version is safe and can be done in-place.

Backwards-compatible bug fixes to Deis are merged into the master branch at any
time after they have :ref:`two approval comments <merge_approval>`.

Patch releases are created as often as needed, based on the priority of one or
more bug fixes that have been merged. If time or severity is crucial, an
individual maintainer can create a patch release without consensus from others.
Patch releases are created from a previous release by cherry-picking specific
bug fixes from the master branch, then applying and pushing the new release tag.


Minor Releases
--------------

A minor release of Deis introduces functionality in a backward-compatible
manner. Upgrading to this version is safe and can be done in-place.

Backwards-compatible functionality changes to Deis are merged into the master
branch after they have :ref:`two approval comments <merge_approval>`, and after
the PR has been assigned to a milestone tracking the minor release.

It is preferable to merge several backwards-compatible functionality changes for
a single minor release.

The Deis project will use GitHub milestones to communicate the content and
timing of planned minor releases. Currently project maintainers meet each week
on Thursday afternoon to assign pull requests to milestones.

The Deis project will release a minor version of the platform on the first
Tuesday of each month. The target day may be shifted to accomodate holidays or
unusual circumstances. If project maintainers agree, an additional minor release
may occur between planned releases. Code freeze and :ref:`further acceptance
tests <features_to_be_tested>` begin on the Thursday before a targeted release.

A minor release may be superceded by a major release.


Major Releases
--------------

A major release of Deis introduces incompatible API changes. Upgrading to this
version may involve a backup and restore process. Custom integrations with Deis
may need to be updated.

Incompatible changes to Deis are merged into the master branch deliberately, by
agreement among maintainers. In addition to
:ref:`two approval comments <merge_approval>`, the pull request must be assigned
to a planning milestone for that release, at which point it can be merged when
release activities and testing begin.

The Deis project will use GitHub milestones to communicate the content and
timing of planned major releases.


.. _`semantic versioning`: http://semver.org/spec/v2.0.0.html
