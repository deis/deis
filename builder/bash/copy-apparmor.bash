set -eo pipefail

# set debug based on envvar
[[ $DEBUG ]] && set -x

main() {
  # cheat
  cp /app/bin/apparmor_parser /sbin/apparmor_parser
  cp /app/bin/auplink /sbin/auplink
}

