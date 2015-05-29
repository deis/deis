:title: Releases
:description: The Deis software release process.

.. _releases:

Release Checklist
=================

This document assists maintainers with creating a new Deis product release.
Please update it to agree with any changes in the process.

Follow these instructions from top to bottom, skipping the sections that do
not apply.


Patch Release
-------------

- Check out the previous release tag

  - ``git checkout vA.B.C``

- Pull in specific bug-fix commits

  - ``git cherry-pick <commit-ish>...``

- Bump the patch version numbers:

  .. code-block:: console

    $ ./contrib/bumpver/bumpver -f A.B.C A.B.D \
        README.md \
        builder/image/Dockerfile \
        builder/image/slugbuilder/Dockerfile \
        builder/image/slugrunner/Dockerfile \
        cache/Dockerfile \
        cache/image/Dockerfile \
        client/deis.py \
        client/setup.py \
        contrib/coreos/user-data.example \
        controller/deis/__init__.py \
        controller/Dockerfile \
        database/Dockerfile \
        deisctl/cmd/cmd.go \
        deisctl/deis-version \
        docs/_includes/_get-the-source.rst \
        docs/installing_deis/install-deisctl.rst \
        docs/installing_deis/install-platform.rst \
        docs/managing_deis/upgrading-deis.rst \
        docs/reference/api-v1.4.rst \
        docs/troubleshooting_deis/index.rst \
        logger/image/Dockerfile \
        logspout/image/Dockerfile \
        publisher/image/Dockerfile \
        registry/Dockerfile \
        router/image/Dockerfile \
        store/base/Dockerfile \
        version/version.go

- Update the CHANGELOG to include all commits since the last release

  - ``./contrib/util/generate-changelog.sh vA.B.C | cat - CHANGELOG.md > tmp && mv tmp CHANGELOG.md``
  - change ``HEAD`` at the top to ``vA.B.D`` (the new release)
  - remove any empty sections and proofread for consistency

- ``git grep A.B.C`` to ensure that no old version strings were missed

- Commit and push the tag

  - ``git commit -a -m 'chore(release): update version to vA.B.D'``
  - ``git tag vA.B.D``
  - ``git push --tags origin vA.B.D``


Major or Minor Release
----------------------

- Move any open issues to the next `deis milestone`_, then close this one
- Check out and update the Deis repo master branch

  - ``git checkout master && git pull``

- Bump the major or minor version numbers

  .. code-block:: console

    $ ./contrib/bumpver/bumpver -f A.B.D-dev A.B.D \
        client/deis.py \
        client/setup.py \
        controller/deis/__init__.py \
        deisctl/deis-version \
        version/version.go

    $ ./contrib/bumpver/bumpver -f A.B.C A.B.D \
        README.md \
        contrib/coreos/user-data.example \
        docs/_includes/_get-the-source.rst \
        docs/installing_deis/install-deisctl.rst \
        docs/installing_deis/install-platform.rst \
        docs/managing_deis/upgrading-deis.rst \
        docs/reference/api-v1.4.rst \
        docs/troubleshooting_deis/index.rst

  - Edit deisctl/cmd/cmd.go and change the default in the RefreshUnits usage string
    (near the bottom of the file) from ``[master]`` to ``[vA.B.D]``.

  - Find and replace "A.B.D-dev" with "A.B.D" in all project Dockerfiles.

- Update the CHANGELOG to include all commits since the last release

  - ``./contrib/util/generate-changelog.sh vA.B.C | cat - CHANGELOG.md > tmp && mv tmp CHANGELOG.md``
  - change ``HEAD`` at the top to ``vA.B.D`` (the new release)
  - remove any empty sections and proofread for consistency

- ``git grep A.B.C`` to ensure that no old version strings were missed

- Commit and push the tag to master

  - ``git commit -a -m 'chore(release): update version to vA.B.D'``
  - ``git push origin master``
  - ``git tag vA.B.D``
  - ``git push --tags origin vA.B.D``


Any Release
-----------

- If this release includes a new component, configure `test-acceptance`_ to publish it to Docker Hub

- Trigger CI jobs manually at https://ci.deis.io/, specifying the new vA.B.D tag

  - build-deis-cli-installer-darwin
  - build-deis-cli-installer-linux
  - build-deisctl-installer-darwin
  - build-deisctl-installer-linux
  - *after* these client jobs finish, trigger test-acceptance

- Publish Deis CLI to pypi.python.org

  - ``pushd client && python setup.py sdist upload && popd``

- Publish docs to pythonhosted.org/deis

  - ``make -C docs clean zipfile``
  - upload docs/docs.zip to the web form at the `Deis pypi`_ page

- Update the installer scripts at `deis/deis.io`_ to reference new version A.B.D

  - https://github.com/deis/deis.io/blob/gh-pages/deis-cli/install.sh
  - https://github.com/deis/deis.io/blob/gh-pages/deisctl/install.sh

- Update published doc versions at ReadTheDocs

  - log in to the https://readthedocs.org/ admin
  - add the current release to the published versions
  - remove the oldest version from the list of published builds
  - rebuild all published versions so their "Versions" index links update

- Update the Homebrew install recipes for ``deis`` and ``deisctl`` with PRs

  - https://github.com/Homebrew/homebrew/blob/master/Library/Formula/deis.rb
    (check for updated python requirements too)
  - https://github.com/Homebrew/homebrew/pull/34967

- Update #deis IRC channel topic to reference new version


Patch Release
-------------

- Bump the version numbers in master to the new release

  .. code-block:: console

    ./contrib/bumpver/bumpver -f A.B.C A.B.D \
      README.md \
      contrib/coreos/user-data.example \
      docs/_includes/_get-the-source.rst \
      docs/installing_deis/install-deisctl.rst \
      docs/installing_deis/install-platform.rst \
      docs/managing_deis/upgrading-deis.rst \
      docs/reference/api-v1.3.rst \
      docs/troubleshooting_deis/index.rst

  - ``git commit -a -m 'chore(release): update version in master to vA.B.D'``
  - ``git push origin master``

- Create `release notes`_ on GitHub

  - copy and paste the newly added CHANGELOG.md section as the body
  - preface with an explanatory paragraph if necessary, for example to reference
    security fixes or point out upgrade details


Major or Minor Release
----------------------

- Edit deisctl/cmd/cmd.go and change the default in the RefreshUnits usage string
  (near the bottom of the file) from ``[vA.B.D]`` to ``[master]``
- Bump the version numbers in master to the next planned with ``-dev``

  .. code-block:: console

    $ ./contrib/bumpver/bumpver -f A.B.D A.B.E-dev \
        client/deis.py \
        client/setup.py \
        controller/deis/__init__.py \
        deisctl/deis-version \
        version/version.go

  - Find and replace "A.B.D" with "A.B.D-dev" in all project Dockerfiles.
  - ``git commit -a -m 'chore(release): update version in master to vA.B.D-dev'``
  - ``git push origin master``

- Create release notes blog post at `deis/deis.io`_ following previous formats
- Create `release notes`_ at GitHub

  - copy and paste from the previous blog post
  - remove Jekyll-specific headers and ``<!-- more -->`` tag

- Update HipChat channel topic to reference the next planned version


.. _`deis milestone`: https://github.com/deis/deis/issues/milestones
.. _`deis open issues`: https://github.com/deis/deis/issues?state=open
.. _`changelog script`: https://github.com/deis/deis/blob/master/contrib/util/generate-changelog.sh
.. _`release notes`: https://github.com/deis/deis/releases
.. _`Deis pypi`:  https://pypi.python.org/pypi/deis/
.. _`deis/deis.io`: https://github.com/deis/deis.io
.. _`test-acceptance`: https://ci.deis.io/job/test-acceptance/configure
