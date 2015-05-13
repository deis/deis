set -eo pipefail

# set debug based on envvar
[[ $DEBUG ]] && set -x

main() {
  HOSTNAME=`hostname`

  if ! etcdctl --no-sync -C $ETCD get ${ETCD_PATH}/monSetupComplete >/dev/null 2>&1 ; then
    echo "store-monitor: Ceph hasn't yet been deployed. Trying to deploy..."
    # let's rock and roll. we need to obtain a lock so we can ensure only one machine is trying to deploy the cluster
    if etcdctl --no-sync -C $ETCD mk ${ETCD_PATH}/monSetupLock $HOSTNAME >/dev/null 2>&1 \
    || [[ `etcdctl --no-sync -C $ETCD get ${ETCD_PATH}/monSetupLock` == "$HOSTNAME" ]] ; then
      echo "store-monitor: obtained the lock to proceed with setting up."

      # Generate administrator key
      ceph-authtool /etc/ceph/ceph.client.admin.keyring --create-keyring --gen-key -n client.admin --set-uid=0 --cap mon 'allow *' --cap osd 'allow *' --cap mds 'allow'

      # Generate the mon. key
      ceph-authtool /etc/ceph/ceph.mon.keyring --create-keyring --gen-key -n mon. --cap mon 'allow *'

      fsid=$(uuidgen)
      etcdctl --no-sync -C $ETCD set ${ETCD_PATH}/fsid ${fsid} >/dev/null

      # Generate initial monitor map
      monmaptool --create --add ${HOSTNAME} ${HOST} --fsid ${fsid} /etc/ceph/monmap

      etcdctl --no-sync -C $ETCD set ${ETCD_PATH}/monKeyring < /etc/ceph/ceph.mon.keyring >/dev/null
      etcdctl --no-sync -C $ETCD set ${ETCD_PATH}/adminKeyring < /etc/ceph/ceph.client.admin.keyring >/dev/null

      # mark setup as complete
      echo "store-monitor: setup complete."
      etcdctl --no-sync -C $ETCD set ${ETCD_PATH}/monSetupComplete youBetcha >/dev/null
    else
      until etcdctl --no-sync -C $ETCD get ${ETCD_PATH}/monSetupComplete >/dev/null 2>&1 ; do
        echo "store-monitor: waiting for another monitor to complete setup..."
        sleep 5
      done
    fi
  fi
}
