:title: Deploying with Heroku Buildpacks on Deis
:description: How to deploy applications using Heroku Buildpacks

Buildpacks
==========

Buildpacks are bundles of detection and configuration scripts which set up containers to
run applications.

Deploy using Buildpacks
-----------------------

For convenience, there are a few buildpacks that are bundled with Deis:

 * `Java buildpack`_
 * `Ruby Buildpack`_
 * `Python Buildpack`_
 * `Nodejs Buildpack`_
 * `Play Buildpack`_
 * `PHP Buildpack`_
 * `Clojure Buildpack`_
 * `Golang Buildpack`_
 * `Scala Buildpack`_
 * `Dart Buildpack`_
 * `Perl Buildpack`_

Deis will cycle through the ``bin/detect`` scripts of each buildpack to match the code you
are pushing.

Adding Custom Buildpacks
------------------------

To add a specific buildpack to your custom Deis cluster, you will need to make the change
in the `builder recipe`_ for the `Deis cookbook`_.

.. note::

    Not all Heroku buildpacks work with Deis due to environmental differences (e.g.
    missing libraries, Heroku-specific environment changes). Test any buildpack before
    using it in production deployments.

Deploying an App with a Custom Buildpack
----------------------------------------

If you want your application to use a specific buildpack that is not included in the list,
you can set the BUILDPACK_URL environment variable for your application before your first
push. For example:

.. code-block:: console

    $ deis config:set BUILDPACK_URL=https://github.com/bacongobbler/heroku-buildpack-jekyll
    === classy-hardtack
    BUILDPACK_URL: https://github.com/bacongobbler/heroku-buildpack-jekyll


.. _`Java buildpack`: https://github.com/heroku/heroku-buildpack-java.git
.. _`Ruby buildpack`: https://github.com/heroku/heroku-buildpack-ruby.git
.. _`Python buildpack`: https://github.com/heroku/heroku-buildpack-python.git
.. _`Nodejs buildpack`: https://github.com/gabrtv/heroku-buildpack-nodejs
.. _`Play buildpack`: https://github.com/heroku/heroku-buildpack-play.git
.. _`PHP buildpack`: https://github.com/CHH/heroku-buildpack-php.git
.. _`Clojure buildpack`: https://github.com/heroku/heroku-buildpack-clojure.git
.. _`Golang buildpack`: https://github.com/kr/heroku-buildpack-go.git
.. _`Scala buildpack`: https://github.com/heroku/heroku-buildpack-scala.git
.. _`Dart buildpack`: https://github.com/igrigorik/heroku-buildpack-dart.git
.. _`Perl buildpack`: https://github.com/miyagawa/heroku-buildpack-perl/tree/carton
.. _`builder recipe`: https://github.com/opdemand/deis-cookbook/blob/master/recipes/builder.rb
.. _`Deis cookbook`: https://github.com/opdemand/deis-cookbook.git
