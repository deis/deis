:title: Releases
:description: Details the Deis release process. Deis releases.
:keywords: deis, release, process, build, tag

.. _releases:

Releases
========

When the maintainers create a **Deis** release, here are the steps involved:


GitHub Issues
-------------

- create next `milestone`_
- roll unfinished issues (if there are any) into next milestone
- close current release milestone


Other Repos
-----------

- tag the opdemand/buildstep repo
- tag the opdemand/gitosis repo


Chef Repo
---------

- checkout release branch
- merge master into release branch locally
- change chef attributes from master to latest tag in deis-cookbook/attributes
    * default.deis.build.revision
    * default.deis.gitosis.revision
    * default.deis.controller.revision
- ``knife cookbook metadata .`` will update metadata.json
-  commit and push the opdemand/deis-cookbook repo
    * ``git commit -a -m 'updated for v0.0.X release'``
    * ``git push origin release``
    * ``git tag v0.0.X``
    * ``git push --tags origin v0.0.X``
- update opscode community cookbook
    * ``cp -pr deis-cookbook /tmp/deis && cd /tmp``
    * ``tar cvfz deis-cookbook-v0.0.X.tar.gz --exclude='deis/.git' --exclude='deis/.vagrant' deis``
    * log in to community.opscode.com and upload tarball
- switch master back to upcoming release
    * ``git checkout master``
    * change cookbook revisions in metadata.rb to *next* version
    * ``git commit -a -m 'switch master to v0.0.Y'``
    * ``git push origin master``


Deis Repo
---------

- merge master into release branch
    * create pull request from master to release branch
    * review and merge on github.com
- Update berksfile with new release
    * ``berks update && berks install``
    * switch from github cookbook to opscode community cookbook
- commit and push the opdemand/deis repo
    * ``git commit -a -m 'updated for v0.0.X release'``
    * ``git tag v0.0.X``
    * ``git push --tags origin v0.0.X``
- switch master to upcoming release
    * ``git checkout master``
    * update __version__ fields in Python packages to *next* version
    * switch from opscode community cookbook back to github cookbook
    * ``berks update && berks install``
    * ``git commit -a -m 'switch master to v0.0.Y'``
    * ``git push origin master``


Client
------
- publish CLI to pypi.python.org
    - ``python setup.py sdist upload``
    - use testpypi.python.org first to ensure there aren't any problems


Docs
----
- create release notes docs
    - follow format of previous `release notes`_
    - summarize all work done
    - what's next and future directions
- publish docs to http://docs.deis.io (deis.readthedocs.org)
- publish docs to pythonhosted.org/deis
    - from the project root, run ``make -C docs clean zipfile``
    - zipfile will be at *docs/docs.zip*
    - log in and use web form at https://pypi.python.org/pypi/deis/
      to upload zipfile


.. _`milestone`: https://github.com/opdemand/deis/issues/milestones
.. _`release notes`: https://github.com/opdemand/deis/releases
