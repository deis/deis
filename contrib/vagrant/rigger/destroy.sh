function _destroy-all-vagrants {
  VMS=$(vagrant global-status | grep deis | awk '{ print $5 }')
  for dir in $VMS; do
    cd ${dir} && vagrant destroy --force
  done
}

rerun_log "Destroying Vagrant cluster..."
_destroy-all-vagrants
