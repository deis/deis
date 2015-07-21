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

function setup-upgrader {
  local upgrader="${1}"

  source "${UPGRADER_DIR}/interface.sh"
  source "${UPGRADER_DIR}/${upgrader}.sh"
}

function source-shared {
  source "${RERUN_MODULE_DIR}/../test-setup.sh"

  # Needs DEIS_TEST_ID
  source "${RERUN_MODULE_DIR}/lib/config.sh"
  source "${RERUN_MODULE_DIR}/lib/clients.sh"
  source "${RERUN_MODULE_DIR}/lib/checks.sh"
  source "${RERUN_MODULE_DIR}/lib/platform.sh"
}

function setup-provider-dependencies {
  _setup-provider-dependencies
}

function destroy-cluster {
  dump-vars

  if [ "${SKIP_CLEANUP}" != true ]; then
    rerun_log "Cleaning up"
    _destroy || true
  fi
}

function create-cluster {
  _create
}

function dump-vars {
  echo
  rerun_log "Dumping useful variables (sourceable: ${TEST_ROOT}/vars)..."

  local output=$(cat <<EOF
TEST_ROOT="${TEST_ROOT}"
DEISCTL_TUNNEL="${DEISCTL_TUNNEL}"
PATH="${PATH}"
DEIS_TEST_SSH_KEY="${DEIS_TEST_SSH_KEY}"
EOF
)

  mkdir -p "${TEST_ROOT}"
  echo "${output}" > "${TEST_ROOT}/vars"
  echo "======= VARIABLES ======="
  echo "${output}"
  echo "========================="
}

source-shared
