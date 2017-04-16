#!/bin/sh
vagrant destroy -f
knife node delete -y deis-controller
knife client delete -y deis-controller
knife data bag delete -y deis-apps
knife data bag delete -y deis-formations
