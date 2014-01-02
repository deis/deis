:title: FAQ
:description: Frequently asked questions about the Deis project. Deis FAQ.
:keywords: deis, PaaS, cloud, faq, custom buildpack

.. _faq:

FAQ
===

- What does the word "deis" mean?

    Deis is an alternative form of dais_, a raised platform for dignified occupancy.

- How do you pronounce "deis?"

    DAY-iss

.. _dais: https://en.wiktionary.org/wiki/dais

- How can I use custom buildpacks with Deis?

    1. Clone the `deis-cookbook`_ repository.
    2. Change the *buildpacks* definition in *recipes/build.rb*.
    3. Upload the changed cookbook to the Chef server
       with ``berks upload --force``.
    4. SSH into your Deis controller and run ``sudo chef-client``.

.. _`deis-cookbook`: https://github.com/opdemand/deis-cookbook.git
