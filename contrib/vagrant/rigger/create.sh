(
  cd ${DEIS_ROOT}
  vagrant up --provider virtualbox
)

export DEISCTL_TUNNEL="${DEISCTL_TUNNEL:-127.0.0.1:2222}"
save-var DEISCTL_TUNNEL
