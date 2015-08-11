function upgrade-deis {
  local from="${1}"
  local to="${2}"

  undeploy-deis

  setup-clients "${to}"

  build-deis "${to}"

  deploy-deis "${to}"
}
