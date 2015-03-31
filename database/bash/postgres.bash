set -eo pipefail

main() {
  # ensure WAL log bucket exists
  envdir /etc/wal-e.d/env /app/bin/create_bucket ${BUCKET_NAME}

  if [[ ! -f /var/lib/postgresql/9.3/main/initialized ]]; then
    echo "database: no existing database found."
    # check if there are any backups -- if so, let's restore
    # we could probably do better than just testing number of lines -- one line is just a heading, meaning no backups
    if [[ `envdir /etc/wal-e.d/env wal-e --terse backup-list 2> /dev/null | wc -l` -gt "1" ]]; then
      echo "database: restoring from backup..."
      rm -rf /var/lib/postgresql/9.3/main
      sudo -u postgres envdir /etc/wal-e.d/env wal-e backup-fetch /var/lib/postgresql/9.3/main LATEST
      chown -R postgres:postgres /var/lib/postgresql/9.3/main
      chmod 0700 /var/lib/postgresql/9.3/main
      echo "restore_command = 'envdir /etc/wal-e.d/env wal-e wal-fetch \"%f\" \"%p\"'" | sudo -u postgres tee /var/lib/postgresql/9.3/main/recovery.conf >/dev/null
    else
      echo "database: initializing a new database..."
    fi
    # either way, we mark the database as initialized
    touch /var/lib/postgresql/9.3/main/initialized
  else
    echo "database: existing data directory found. Starting postgres..."
  fi

  # perform a one-time reload to populate database entries
  /usr/local/bin/reload
}