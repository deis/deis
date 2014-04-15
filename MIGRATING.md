Migrating to Deis 0.8.0
=======================

If you are updating from Deis version 0.7.0 or earlier, there are
several big changes you should know about.

If you need to use Deis with Chef integration, on Ubuntu 12.04 LTS, or
on DigitalOcean, you should use the
[v0.7.0 release](https://github.com/opdemand/deis/tree/v0.7.0) of Deis.

Upgrading
---------
During the last several months of pre-release development of Deis, we
have not been able to provide a smooth upgrade process between every
release. Unfortunately, that is also the case for the upcoming release,
v0.8.0. The reasons why a simple in-place upgrade isn't possible for
this release are made plain below.

However, the Deis core developers believe the current release is stable
architecturally--we don't anticipate major changes to components,
interfaces, or packaging at this point. And we will soon implement a
reliable upgrade process that will enable you to stay with Deis changes
from here on out by simple in-place upgrades.

So thanks for your patience. We could not have gotten here without you,
the Deis community. Now we're here to support you, so install Deis 0.8.0
and come on along with us.

Chef not Required
-----------------
Versions of Deis previous to v0.8.0 relied on Chef for provisioning
a controller and nodes, for tracking users and apps, and for deploying
applications. Chef is now no longer a requirement to use Deis.

Now anywhere you can run CoreOS, you can run Deis. Provisioning a
contoller and nodes can be done through a variety of IT-friendly methods
available for everything from Rackspace Cloud to Vagrant to PXE hardware
boot--see the CoreOS documentation for details.

Users and apps are tracked in Deis' controller and database, and apps
are deployed across a cluster of machines by the `fleet` distributed
init system.

CoreOS replaces Ubuntu
----------------------
Versions of Deis previous to v0.8.0 supported installation on Ubuntu
12.04 LTS. Deis now runs on CoreOS, a lean distribution targeted at
massive server deployments.

No More Formations or Layers
----------------------------
Versions of Deis previous to v0.8.0 organized groups of machines into
"formations" with internal "layers" used to scale nodes. Deis now
simplifies both concepts into that of a "cluster": a set of machines
available for smart process scheduling.

No More Providers
-----------------
Deis now incorporates new machines into a cluster through whatever IT
process is appropriate; there is no more `deis nodes:scale`. Individual
support modules for cloud providers such as Rackspace and Amazon EC2
are no longer part of Deis' core code, and there is no need to import
credentials through the Deis CLI.

Amazon EC2, Rackspace and other cloud providers, as well as bare metal
deployments, continue to be the focus of support and development for
Deis. But Deis no longer needs vendor-specific code in order to grow
a cluster or scale an application.

Support scripts and documentation for Amazon EC2, Rackspace, and other
cloud providers are always available in the
[contrib](contrib/) directory.

DigitalOcean Support
--------------------
DigitalOcean does not yet support deploying CoreOS on droplets. Until
there is a CoreOS solution, Deis cannot support clusters
on DigitalOcean.

If you use DigitalOcean, please
[show your support](http://digitalocean.uservoice.com/forums/136585-digital-ocean/suggestions/4250154-suport-coreos-as-a-deployment-platform)
for CoreOS and help us to support Deis on DO.

deis/deis in GitHub
-------------------
The https://github.com/opdemand/deis will soon move under the banner of
*The Deis Project* at https://github.com/deis, so getting the source
code will be:

```console
$ git clone https://github.com/deis/deis.git
```

Nothing else about the project changes; we just think it will be easier
to find Deis there, next to related projects such as deis/tester,
deis/base, and deis/slugbuilder.
