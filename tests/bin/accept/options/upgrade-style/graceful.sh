function upgrade-deis {
  local from="${1}"
  local to="${2}"

  deisctl upgrade-prep

  healthcheck

  setup-clients "${to}"

  build-deis "${to}"

  deisctl upgrade-takeover
}
