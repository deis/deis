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

Deis Repo
---------

- ``bundle install``
- Update berksfile with new release
    * ``berks update && berks install``
    * switch from github cookbook to opscode community cookbook
- tag the opdemand/deis-cookbook repo
	* ``git status && git add . && git commit -m 'updating for 0.0.X release'``
	* ``git tag v0.0.X``
	* ``git push origin master``
	* ``git push --tags``
- update __version__ fields in Python packages to *next* versionÏ
- switch from opscode community cookbook back to github cookbook

Chef Repo
---------

- change chef attributes from master to latest tag in deis-cookbook/attributes
	- default.deis.build.revision
	- default.deis.gitosis.revision
	- default.deis.controller.revision
- ``knife cookbook metadata .`` will update metadata.json
- tag the opdemand/deis-cookbook repo
	* ``git commit -a -m 'prep for 0.0.X release'``
	* ``git tag v0.0.X``
	* ``git push origin master``
	* ``git push --tags``
- ``cp -pr deis-cookbook /tmp/deis && cd /tmp``
- ``tar cvfz deis-cookbook-v0.0.6.tar.gz --exclude='deis/.git' --exclude='deis/.vagrant' deis``
- log in to community.opscode.com and upload tarball
- change gitosis, build, controller from latest back to master tag
- change cookbook revisions in metadata.rb to *next* version
- git commit and push post-tag dev versions

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
