#!/bin/sh
#
# Allows git to exec SSH but bypass auth warnings.
# To use, export the environment variable GIT_SSH as the location of this script,
# then run git commands as usual:
# $ export GIT_SSH=$HOME/bin/git-ssh-nokeycheck.sh

# We want the $@ to expand client side
# shellcheck disable=SC2029
ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null "$@"
