:title: Releases
:description: The Deis software release process.

.. _releases:

Release Checklist
=================

These instructions are to assist the Deis maintainers with creating a new Deis
product release. Please keep this document up-to-date with any changes in this process.

deis repo
---------
- If this release was managed as a milestone in GitHub:
    * Create the next `deis milestone`_
    * Move any `deis open issues`_ from the current release to the next milestone
    * Close the current `deis milestone`_
- Create a branch for the release PR: ``git checkout -b release-X.Y.Z``
- Update CHANGELOG.md using the `changelog script`_
    * ``./contrib/util/generate-changelog.sh vU.V.W | cat - CHANGELOG.md > tmp && mv tmp CHANGELOG.md``
      substituting the previous release for vU.V.W.
    * proofread the new CHANGELOG.md to ensure it was generated correctly and edit ``HEAD`` at the top
      to vX.Y.Z (the current release)
- Update version strings with the ``bumpver`` tool:

  .. code-block:: console

    $ ./contrib/bumpver/bumpver X.Y.Z \
        version/version.go \
        README.md \
        client/deis.py \
        client/setup.py \
        contrib/coreos/user-data.example \
        controller/deis/__init__.py \
        deisctl/deis-version \
        deisctl/deisctl.go \
        docs/contributing/test_plan.rst \
        docs/installing_deis/install-deisctl.rst \
        docs/installing_deis/install-platform.rst \
        docs/managing_deis/upgrading-deis.rst

- Edit deisctl/cmd/cmd.go and change the default in the RefreshUnits usage string
  (near the bottom of the file) from ``[master]`` to ``[vX.Y.Z]``.
- Examine the output of ``git grep vU.V.W`` to ensure that no old version strings
  were missed
- Commit and push the deis/deis release and tag
    * ``git commit -a -m 'chore(release): update version to vX.Y.Z'``
    * ``git push origin release-X.Y.Z``
- When the PR is approved and merged, tag it in master
    * ``git checkout master && git pull``
    * ``git tag vX.Y.Z``
    * ``git push --tags origin vX.Y.Z``
- Trigger all deis-cli and deisctl builder jobs at ci.deis.io. When they finish, verify that
  the current binary installers are publicly downloadable from the opdemand S3 bucket.
- Trigger the test-master job, supplying vX.Y.Z as the version
- When test-master completes, double-check images at Docker Hub to verify tags are published
- Publish CLI to pypi.python.org
    - ``cd client && python setup.py sdist upload``
    - use testpypi.python.org first to ensure there aren't any problems

deis.io repo
------------
- Update deis.io installer scripts to point to new versions by default
    * update https://github.com/deis/deis.io/blob/gh-pages/deis-cli/install.sh
    * update https://github.com/deis/deis.io/blob/gh-pages/deisctl/install.sh

Documentation
-------------
- (CHANGELOG.md files were regenerated and committed above.)
- Docs are automatically published to http://docs.deis.io (the preferred alias
  for deis.readthedocs.org)
- Log in to the http://deis.readthedocs.org admin
    * add the current release to the list of published builds
    * remove the oldest release from the list of published builds
    * rebuild all published versions so their "Versions" index links
      are updated
- Publish docs to pythonhosted.org/deis
    * from the project root, run ``make -C docs clean zipfile``
    * the zipfile will be at **docs/docs.zip**
    * log in to http://pypi.python.org/ and use the web form at the
      `Deis Pypi`_ page to upload the zipfile
- Check documentation for deis/* projects at the `Docker Hub`_
    * click "Settings" for each project (deis/controller, deis/cache, etc.)
    * paste the contents of each README.md into the "long description" field if
      there are discrepencies
- For a milestone release, create release notes docs
    * follow the format of previous `release notes`_
    * summarize all work done since the previous release
    * visit all deis/* project issues to make sure we don't
      miss any contributors for the "Community Shout-Outs" section
    * include "what's next" and "future directions" sections
    * add Markdown version of release notes to `deis/deis.io`_ website project
- For a patch release, paste the new CHANGELOG.md section as GitHub release notes

Post-Release
------------
- Update the #deis IRC channel topic to reference the newly released version
- For a milestone release, update HipChat channel topics to reference the
  next planned version
- Create a branch for the post-release PR: ``git checkout -b release-X.Y.Z+git``
- Update version strings to vX.Y.Z+git with the ``bumpver`` tool:

  .. code-block:: console

    $ ./contrib/bumpver/bumpver X.Y.Z+git \
        version/version.go \
        client/deis.py \
        deisctl/deis-version \
        deisctl/deisctl.go \
        controller/deis/__init__.py \
        README.md

- Edit deisctl/cmd/cmd.go and change the default in the RefreshUnits usage string
  (near the bottom of the file) from ``[vX.Y.Z]`` to ``[master]``.
- Create a pull request for vX.Y.Z+git
    * ``git commit -a -m 'chore(release): update version in master to vX.Y.Z+git'``
- Ensure that this PR is merged before others are allowed to be merged!


.. _`deis milestone`: https://github.com/deis/deis/issues/milestones
.. _`deis open issues`: https://github.com/deis/deis/issues?state=open
.. _`changelog script`: https://github.com/deis/deis/blob/master/contrib/util/generate-changelog.sh
.. _`release notes`: https://github.com/deis/deis/releases
.. _`aws-eng S3 bucket`: https://s3-us-west-2.amazonaws.com/opdemand/
.. _`Deis Pypi`:  https://pypi.python.org/pypi/deis/
.. _`Docker Hub`: https://hub.docker.com/
.. _`deis/deis.io`: https://github.com/deis/deis.io
