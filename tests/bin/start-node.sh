#!/usr/bin/env bash
#
# Connects this machine as a Jenkins node to https://ci.deis.io/
# Set NODE_NAME and NODE_SECRET to the values provided by Jenkins in the "Manage Nodes"
# administrative interface.

wget -N https://ci.deis.io/jnlpJars/slave.jar
java -jar slave.jar -jnlpUrl https://ci.deis.io/computer/${NODE_NAME?}/slave-agent.jnlp -secret ${NODE_SECRET?} &
