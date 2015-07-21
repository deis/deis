# Shell functions for the tests module.
#/ usage: source RERUN_MODULE_DIR/lib/functions.sh command
#

# Read rerun's public functions
. $RERUN || {
    echo >&2 "ERROR: Failed sourcing rerun function library: \"$RERUN\""
    return 1
}

# Check usage. Argument should be command name.
[[ $# = 1 ]] || rerun_option_usage

# Source the option parser script.
#
if [[ -r $RERUN_MODULE_DIR/commands/$1/options.sh ]] 
then
    . $RERUN_MODULE_DIR/commands/$1/options.sh || {
        rerun_die "Failed loading options parser."
    }
fi

# - - -
# Your functions declared here.
# - - -

function not-implemented {
  rerun_log "No implementation of ${FUNCNAME[1]} in ${PROVIDER}"
}

function setup-provider {
# This loads the provider's implementations of the provider interface
  local provider="${1}"

  source "${PROVIDER_DIR}/interface.sh"
  source "${PROVIDER_DIR}/${provider}.sh"
}

function source-shared {
  source "${RERUN_MODULE_DIR}/lib/config.sh"
  source "${RERUN_MODULE_DIR}/bin/test-setup.sh"
  source "${RERUN_MODULE_DIR}/lib/clients.sh"
  source "${RERUN_MODULE_DIR}/lib/checks.sh"
}

function setup-provider-dependencies {
  _setup-provider-dependencies ${@}
}

function destroy-cluster {
  if [ "${SKIP_CLEANUP}" != true ]; then
    rerun_log "Cleaning up"
    _destroy
  fi
}

function create-cluster {
  _create ${@}
}

function deploy-deis {
  local version="${1}"

  check-etcd-alive

  deisctl config platform set domain="${DEIS_TEST_DOMAIN}"
  deisctl config platform set sshPrivateKey="${DEIS_TEST_SSH_KEY}"

  if is-released-version "${version}"; then
    deisctl config platform set version="v${version}"
  else
    export DEV_REGISTRY="registry.hub.docker.com"
    export IMAGE_PREFIX="deisci/"

    make build dev-release
  fi

  deisctl install platform
  deisctl start platform
}

function undeploy-deis {
  deisctl stop platform
  deisctl uninstall platform
}

source-shared
