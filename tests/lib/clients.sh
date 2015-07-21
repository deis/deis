function setup-clients {
  local version="${1}"

  rerun_log "Installing clients (version ${version})"

  setup-deis-client "${version}"
  setup-deisctl-client "${version}"
}

function is-released-version {
  [[ ${1} =~ ^([0-9]+\.){0,2}[0-9]$ ]] && return 0
}

function download-client {
  local client="${1}"
  local version="${2}"
  local dir="${3}"

  pushd "${dir}"
    curl -sSL "http://deis.io/${client}/install.sh" | sh -s "${version}"
  popd
}

function setup-deis-client {
# install deis CLI from http://deis.io/ website
  local version="${1}"
  
  local client_dir="${TEST_ROOT}/${version}/deis-cli"
  mkdir -p "${client_dir}"

  if is-released-version "${version}" ; then
    rerun_log "Installing deis-cli (${version}) to ${client_dir}..."
    
    download-client "deis-cli" "${version}" "${client_dir}"
    
    export PATH="${client_dir}:${PATH}"
  else
    rerun_log "Building deis-cli locally..."
    pushd "${DEIS_ROOT}"
      git fetch
      git checkout "${version}"
      make -C client build
      export PATH="${client_dir}/client/dist:${PATH}"
    popd
  fi
}

function setup-deisctl-client {
# install deisctl from http://deis.io/ website or from repository
  local version="${1}"

  unset DEISCTL_UNITS

  local client_dir="${TEST_ROOT}/${version}/deisctl"
  mkdir -p "${client_dir}"

  if is-released-version "${version}" ; then
    rerun_log "Installing deisctl (${version}) to ${client_dir}..."

    download-client "deisctl" "${version}" "${client_dir}"
  else
    rerun_log "Building deisctl locally..."

    pushd "${DEIS_ROOT}"
      git fetch
      git checkout "${version}"
      make -C deisctl build
    popd

    export DEISCTL_UNITS="${DEIS_ROOT}/deisctl/units"
    export PATH="${client_dir}/deisctl:${PATH}"
    deisctl refresh-units
  fi
}
