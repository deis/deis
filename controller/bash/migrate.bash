set -eo pipefail

# set debug based on envvar
[[ $DEBUG ]] && set -x

main() {
  cd /app

  mkdir -p /data/logs
  chmod 777 /data/logs

  # run an idempotent database migration
  sudo -E -u deis /app/manage.py syncdb --migrate --noinput
}
