function build-deis {
  local version="${1}"

  if is-released-version "${version}"; then
    deisctl config platform set version="v${version}"
  else
    make build dev-release
  fi
}

function deploy-deis {
  local version="${1}"

  check-etcd-alive

  deisctl config platform set domain="${DEIS_TEST_DOMAIN}"
  deisctl config platform set sshPrivateKey="${DEIS_TEST_SSH_KEY}"

  build-deis "${version}"

  deisctl install platform
  deisctl start platform

  _check-cluster
}

function undeploy-deis {
  deisctl stop platform
  deisctl uninstall platform
}
