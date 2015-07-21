#!/usr/bin/env roundup
#
#/ usage:  rerun stubbs:test -m accept -p clients [--answers <>]
#

# Helpers
# -------
[[ -f ./functions.sh ]] && . ./functions.sh

# The Plan
# --------
describe "clients"

source ../lib/clients.sh

TEST_ROOT="$(mktemp -d /tmp/roundup-test.XXX)"

function cleanup {
  rm -r "${TEST_ROOT}"
}

trap cleanup EXIT

it_sets_up_released_deisctl() {

  function download-client {
    [ ${1} == "deisctl" ] &&
    [ ${2} == "1.8.0" ]
  }

  setup-deisctl-client "1.8.0"
}

it_sets_up_build_of_deisctl() {

  function git {
    if [ "${1}" == "fetch" ]; then
      return 0
    elif [ "${1}" == "checkout" ] && [ "${2}" == "dev" ]; then
      return 0
    else
      return 1
    fi
  }

  function make {
    [ "$*" == "-C deisctl build" ]
  }

  function deisctl {
    [ "${1}" == "refresh-units" ]
  }

  setup-deisctl-client "dev"

}

it_sets_up_released_deiscli() {

  function download-client {
    [ ${1} == "deis-cli" ] &&
    [ ${2} == "1.8.0" ]
  }

  setup-deis-client "1.8.0"
}

it_sets_up_build_of_deiscli() {

  function git {
    if [ "${1}" == "fetch" ]; then
      return 0
    elif [ "${1}" == "checkout" ] && [ "${2}" == "dev" ]; then
      return 0
    else
      return 1
    fi
  }

  function make {
    [ "$*" == "-C client build" ]
  }

  setup-deis-client "dev"
}

it_identifies_released_versions() {
  is-released-version "1.8.0"
}

it_identifies_unreleased_versions() {
  ! is-released-version "dev"
}
