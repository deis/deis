set -eo pipefail

main() {
  if [[ -f /var/lib/postgresql/9.3/main/recovery.conf ]] ; then
    echo "database: database is currently recovering from a backup. Will try again next time..."
  else
    # perform a backup
    sudo -u postgres envdir /etc/wal-e.d/env wal-e backup-push /var/lib/postgresql/9.3/main
    # only retain the latest BACKUPS_TO_RETAIN backups
    sudo -u postgres envdir /etc/wal-e.d/env wal-e delete --confirm retain ${BACKUPS_TO_RETAIN}
  fi
}