set -eo pipefail

main() {
  # initialize database if one doesn't already exist
  # for example, in the case of a data container
  if [[ ! -d /var/lib/postgresql/9.3/main ]]; then
    chown -R postgres:postgres /var/lib/postgresql
    sudo -u postgres /usr/lib/postgresql/9.3/bin/initdb -D /var/lib/postgresql/9.3/main
  fi
}