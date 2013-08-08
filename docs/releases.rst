:title: Releases
:description: The Deis Release Process
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


Chef Repo
---------

- change cookbook revisions
- change chef attributes deis-cookbook/attributes
	- default.deis.build.revision
	- default.deis.controller.revision
- change chef metadata.rb
- upload cookbook to Chef
	* ``berks update && berks install && berks upload --force``
- tag the opdemand/deis-cookbook repo
	* ``git commit -a -m 'prep for 0.0.X release'``
	* ``git tag v0.0.X``
	* ``git push --tags``


Deis Repo
---------

- ``bundle install``
- Update berksfile with new release
	* ``berks update && berks install && berks upload --force``
- update __version__ fields in Python packages
- tag the opdemand/deis-cookbook repo
	* ``git status && git add . && git commit -m 'updating for 0.0.X release'``
	* ``git tag v0.0.X``
	* ``git push --tags``
- tag the opdemand/buildstep repo
- tag the opdemand/gitosis repo

Client
------

- publish CLI to pip
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
