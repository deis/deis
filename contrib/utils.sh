#!/usr/bin/env bash

function echo_yellow {
  echo -e "\033[0;33m$1\033[0m"
}

function echo_red {
  echo -e "\033[0;31m$1\033[0m"
}

function echo_green {
  echo -e "\033[0;32m$1\033[0m"
}

export COREOS_CHANNEL=${COREOS_CHANNEL:-stable}
export COREOS_VERSION=${COREOS_VERSION:-899.13.0}
export DEIS_RELEASE=1.13.0
