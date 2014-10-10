:title: Releases
:description: The Deis software release process.

.. _releases:

Release Checklist
=================

These instructions are to assist the Deis maintainers with creating a new Deis
product release. Please keep this document up-to-date with any changes in this process.

deis repo
---------
- Create the next `deis milestone`_
- Move any `deis open issues`_ from the current release to the next milestone
- Close the current `deis milestone`_
- Update CHANGELOG.md using the `changelog script`_
    * ``./contrib/util/generate-changelog.sh vU.V.W vX.Y.Z | cat - CHANGELOG.md > tmp && mv tmp CHANGELOG.md``
      substituting the previous release for vU.V.W and the current one for vX.Y.Z.
    * proofread the new CHANGELOG.md to ensure it was generated correctly
    * ``git add CHANGELOG.md && git commit -m "docs(CHANGELOG): update for v.X.Y.Z"``
- Update version strings with the ``bumpver`` tool:

  .. code-block:: console

    $ ./contrib/bumpver/bumpver X.Y.Z \
        version/version.go \
        client/deis.py \
        client/setup.py \
        deisctl/deis-version \
        deisctl/deisctl.go \
        deisctl/README.md \
        contrib/coreos/user-data \
        controller/deis/__init__.py \
        README.md \
        tests/bin/test-latest.sh

- Commit and push the deis/deis release and tag
    * ``git commit -a -m 'chore(release): update version to vX.Y.Z'``
    * ``git push origin master``
    * ``git tag vX.Y.Z``
    * ``git push --tags origin vX.Y.Z``
- Publish CLI to pypi.python.org
    - ``cd client && python setup.py sdist upload``
    - use testpypi.python.org first to ensure there aren't any problems
- Continuous delivery jobs at ci.deis.io will update the deis CLI. Double-check that the
  current binary installers are publicly downloadable from the opdemand S3 bucket.

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
- Create release notes docs
    * follow the format of previous `release notes`_
    * summarize all work done since the previous release
    * visit all deis/* project issues to make sure we don't
      miss any contributors for the "Community Shout-Outs" section
    * include "what's next" and "future directions" sections
    * add Markdown version of release notes to `deis/deis.io`_ website project


.. _`deis milestone`: https://github.com/deis/deis/issues/milestones
.. _`deis open issues`: https://github.com/deis/deis/issues?state=open
.. _`changelog script`: https://github.com/deis/deis/blob/master/contrib/util/generate-changelog.sh
.. _`release notes`: https://github.com/deis/deis/releases
.. _`aws-eng S3 bucket`: https://s3-us-west-2.amazonaws.com/opdemand/
.. _`Deis Pypi`:  https://pypi.python.org/pypi/deis/
.. _`Docker Hub`: https://hub.docker.com/
.. _`deis/deis.io`: https://github.com/deis/deis.io
