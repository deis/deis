set -eo pipefail

# set debug based on envvar
[[ $DEBUG ]] && set -x

main() {
  # this files are created by confd.
  # to avoid an initial error we create empty files
  echo '#!/bin/sh' > /usr/local/bin/reload
  chmod 0755 /usr/local/bin/reload

  touch /etc/postgresql/9.3/main/pg_hba.conf
  touch /etc/postgresql/9.3/main/postgresql.conf
}
