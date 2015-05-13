set -eo pipefail

# set debug based on envvar
[[ $DEBUG ]] && set -x

main() {
  DOCKER="/usr/bin/docker"
  DRIVER_OVERRIDE=$(cat /etc/docker.env)
  # spawn a docker daemon to build cedarish
  rm -rf /var/run/docker.sock
  rm -rf /var/run/docker.pid

  sudo $DOCKER -d --bip=172.19.42.1/16 $DRIVER_OVERRIDE --insecure-registry 10.0.0.0/8 --insecure-registry 172.16.0.0/12 --insecure-registry 192.168.0.0/16 --insecure-registry 100.64.0.0/10 &
  DOCKER_PID=$!

  # wait until the daemon is available
  sleep 10

  # build required images
  $DOCKER build -t deis/slugbuilder /app/slugbuilder/
  $DOCKER build -t deis/slugrunner /app/slugrunner/

  # cleanup.
  kill -SIGTERM $DOCKER_PID
  wait $DOCKER_PID
  rm -rf /var/run/docker.sock
  rm -rf /var/run/docker.pid

  echo "cedarish build finished"
}

