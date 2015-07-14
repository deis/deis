set -eo pipefail

# set debug based on envvar
[[ $DEBUG ]] && set -x

main() {
  echo "$HOST $(hostname)" >> /etc/hosts
}
