:title: Releases
:description: Details the Deis release process. Deis releases.
:keywords: deis, release, process, build, tag

.. _releases:

Release Checklist
=================

These instructions are to assist the Deis maintainers with creating a new Deis
product release.

.. image:: http://upload.wikimedia.org/wikipedia/commons/3/37/Grace_Hopper_and_UNIVAC.jpg
  :width: 220
  :height: 193
  :align: right
  :alt: Grace Hopper and UNIVAC, from the Wikimedia Commons

Please keep this document up-to-date with any changes in this process.

opdemand/deis-cookbook Chef Repo
--------------------------------
- One-time setup:
    * ``gem install stove --pre``
    * Configure your ``~/.stove`` file as outlined in `Installing Stove`_
- Create the next `deis-cookbook milestone`_
- Move any `deis-cookbook open issues`_ from the current release to the
  next milestone
- Close the current `deis-cookbook milestone`_
- Recreate CHANGELOG.md in the root of the project
    * ``npm install github-changes``
    * ``github-changes -o opdemand -r deis-cookbook -n vX.Y.Z``
    * proofread the new CHANGELOG.md to ensure it was generated correctly
    * ``git add CHANGELOG.md && git commit -m "Updated CHANGELOG.md."``
- Merge git master into release branch locally
    * ``git checkout release && git merge master``
- Run ``knife cookbook metadata .`` to update **metadata.json**. **DOUBLE-CHECK
  the generated file after this step--this may be broken currently.**
- Commit, push, and tag the opdemand/deis-cookbook release
    * ``git commit -a -m 'Updated for vX.Y.Z release.'``
    * ``git push origin release``
    * ``git tag vX.Y.Z``
    * ``git push --tags origin vX.Y.Z``
- Update the deis Opscode Community cookbook
    * ``bake vX.Y.Z --no-bump --no-changelog --no-dev --no-git --no-github --no-jira``
- switch master to upcoming release
    * ``git checkout master``
    * change **version** string in metadata.rb to *next* version
    * ``git commit -a -m 'Switch master to vA.B.C.'`` (**next** version)
    * ``git push origin master``

opdemand/deis Server Repo
-------------------------
- Create the next `deis milestone`_
- Move any `deis open issues`_ from the current release to the
  next milestone
- Close the current `deis milestone`_
- Recreate CHANGELOG.md in the root of the project
    * ``github-changes -o opdemand -r deis -n vX.Y.Z``
    * proofread the new CHANGELOG.md to ensure it was generated correctly
    * ``git add CHANGELOG.md && git commit -m "Updated CHANGELOG.md."``
- Merge git master into release branch locally
    * ``git checkout release && git merge master``
- Update Berksfile with new release
    * edit **Berksfile** and ensure it points to the Opscode Community cookbook
      for Deis
    * ``berks update && berks install`` to update **Berksfile.lock**
- Commit and push the opdemand/deis release and tag
    * ``git commit -a -m 'Updated for vX.Y.Z release.'``
    * ``git push origin release``
    * ``git tag vX.Y.Z``
    * ``git push --tags origin vX.Y.Z``
- Publish CLI to pypi.python.org
    - ``cd client && python setup.py sdist upload``
    - use testpypi.python.org first to ensure there aren't any problems
- Create CLI binaries for Windows, Mac OS X, Debian
    - ``pip install pyinstaller && make -C controller client``
    - build **deis-osx-X.Y.Z.tgz** on Mac OS X 10.8 for all Macs (10.9 uses
      LLVM, which makes our binary crash on earlier OS versions)
    - build **deis-win64-X.Y.Z.zip** on Windows 7 64-bit
    - build **deis-deb-wheezy-X.Y.Z.tgz** on Debian Wheezy
      (see https://github.com/opdemand/deis/issues/504)
    - upload all binaries to the `aws-eng S3 bucket`_ and set each as
      publically downloadable
- Switch master to upcoming release
    * ``git checkout master``
    * update __version__ fields in Python packages to *next* version
    * switch from opscode community cookbook back to github cookbook
    * ``berks update && berks install`` to update Berksfile.lock
    * ``git commit -a -m 'Switch master to vA.B.C.'`` (**next** version)
    * ``git push origin master``

Documentation
-------------
- (CHANGELOG.md files were regenerated and committed above.)
- Docs are automatically published to http://docs.deis.io (the preferred alias
  for deis.readthedocs.org)
- Log in to the http://deis.readthedocs.org admin
    * add the current release to the list of published builds
    * rebuild all published versions so their "Versions" index links
      are updated
- Publish docs to pythonhosted.org/deis
    * from the project root, run ``make -C docs clean zipfile``
    * the zipfile will be at **docs/docs.zip**
    * log in to http://pypi.python.org/ and use the web form at the
      `Deis Pypi`_ page to upload the zipfile
- Check documentation for deis/* projects at the `Docker Index`_
    * click "Settings" for each project (deis/controller, deis/cache, etc.)
    * paste the contents of each README.md into the "long description" field if
      there are discrepencies. (These don't automatically sync up after the
      Trusted Build is first created.)
- Create release notes docs
    * follow the format of previous `release notes`_
    * summarize all work done since the previous release
    * visit all opdemand/* and deis/* project issues to make sure we don't
      miss any contributors for the "Community Shout-Outs" section
    * include "what's next" and "future directions" sections


.. _`deis-cookbook milestone`: https://github.com/opdemand/deis-cookbook/issues/milestones
.. _`deis-cookbook open issues`: https://github.com/opdemand/deis-cookbook/issues?state=open
.. _`Opscode Community`: http://community.opscode.com/cookbooks/deis/versions/new
.. _`deis milestone`: https://github.com/opdemand/deis/issues/milestones
.. _`deis open issues`: https://github.com/opdemand/deis/issues?state=open
.. _`release notes`: https://github.com/opdemand/deis/releases
.. _`aws-eng S3 bucket`: https://s3-us-west-2.amazonaws.com/opdemand/
.. _`Deis Pypi`:  https://pypi.python.org/pypi/deis/
.. _`Docker Index`: https://index.docker.io/
.. _`Installing Stove`: https://github.com/sethvargo/stove#installation
