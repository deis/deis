function _setup-provider-dependencies {
  :
}

function _create {
  rerun_log "Creating Vagrant cluster..."
  vagrant up --provider virtualbox

  export DEISCTL_TUNNEL="127.0.0.1:2222"
}

function _destroy {
  rerun_log "Destroying Vagrant cluster..."
  ${BIN_DIR}/destroy-all-vagrants.sh
}
