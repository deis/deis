:title: Releases
:description: The Deis Release Process
:keywords: deis, release, process, build, tag

.. _releases:

Releases
========

When we create a **deis** release, here are the steps involved:

GitHub Issues
-------------

- roll unfinished issues into next milestone
- close release milestone


Chef Repo
---------

- change cookbook revisions
- change chef attributes deis-cookbook/attributes
	- default.deis.build.revision
	- default.deis.controller.revision
- change chef metadata.rb
- upload cookbook to Chef
- berks update && berks install && berks upload --force
- git commit -a -m 'prep for 0.0.X release'
- tag the opdemand/deis-cookbook repo
- git tag v0.0.X
- git push --tags


Deis Repo
---------

- bundle install
- Update berksfile with new release
- berks update && berks install && berks upload --force
- update __version__ fields in packages
- git status && git add . && git commit -m 'updating for 0.0.X release'
- git tag v0.0.X
- git push --tags
- tag the opdemand/deis repo
- tag the opdemand/buildstep repo
- tag the opdemand/gitosis repo

Client
------

- publish CLI to pip
	- python setup.py sdist upload


Docs
----
- create release notes docs
	- summary of features
	- what's next? section
- publish docs to docs.deis.io (TODO)
- publish docs to pythonhosted.org/deis
    - use web form at http://pypi.python.org/pypi?%3Aaction=pkg_edit&name=deis
