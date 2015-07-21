function check-etcd-alive {
  rerun_log "Waiting for etcd/fleet at ${DEISCTL_TUNNEL}"

  # wait for etcd up to 5 minutes
  WAIT_TIME=1
  until deisctl --request-timeout=1 list >/dev/null 2>&1; do
     (( WAIT_TIME += 1 ))
     if [ ${WAIT_TIME} -gt 300 ]; then
      log_phase "Timeout waiting for etcd/fleet"
      # run deisctl one last time without eating the error, so we can see what's up
      deisctl --request-timeout=1 list
      exit 1;
    fi
  done

  rerun_log "etcd available after ${WAIT_TIME} seconds"
}
