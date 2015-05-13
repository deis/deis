set -eo pipefail

# set debug based on envvar
[[ $DEBUG ]] && set -x

main() {
  export PATH=$PATH:/usr/local/bin

  # ensure WAL log bucket exists
  envdir /etc/wal-e.d/env /app/bin/create_bucket ${BUCKET_NAME}
  INIT_ID=$(etcdctl --no-sync -C $ETCD get $ETCD_PATH/initId 2> /dev/null || echo none)
  echo "database: expecting initialization id: $INIT_ID"

  etcdctl --no-sync -C $ETCD set $ETCD_PATH/initialBackup 0 >/dev/null
  if [[ "$(cat /var/lib/postgresql/9.3/main/initialized 2> /dev/null)" != "$INIT_ID" ]]; then
    echo "database: no existing database found or it is outdated."
    # check if there are any backups -- if so, let's restore
    # we could probably do better than just testing number of lines -- one line is just a heading, meaning no backups
    if [[ `envdir /etc/wal-e.d/env wal-e --terse backup-list | wc -l` -gt "1" ]]; then
      echo "database: restoring from backup..."
      rm -rf /var/lib/postgresql/9.3/main
      sudo -u postgres envdir /etc/wal-e.d/env wal-e backup-fetch /var/lib/postgresql/9.3/main LATEST
      chown -R postgres:postgres /var/lib/postgresql/9.3/main
      chmod 0700 /var/lib/postgresql/9.3/main
      echo "restore_command = 'envdir /etc/wal-e.d/env wal-e wal-fetch \"%f\" \"%p\"'" | sudo -u postgres tee /var/lib/postgresql/9.3/main/recovery.conf >/dev/null
    else
      echo "database: no backups found. Initializing a new database..."
      etcdctl --no-sync -C $ETCD set $ETCD_PATH/initialBackup 1 >/dev/null
    fi
    # either way, we mark the database as initialized
    INIT_ID=$(cat /proc/sys/kernel/random/uuid)
    echo $INIT_ID > /var/lib/postgresql/9.3/main/initialized
    etcdctl --no-sync -C $ETCD set $ETCD_PATH/initId $INIT_ID >/dev/null
  else
    echo "database: existing data directory found. Starting postgres..."
  fi

  # perform a one-time reload to populate database entries
  /usr/local/bin/reload
}