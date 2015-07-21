set -eo pipefail

# set debug based on envvar
[[ $DEBUG ]] && set -x

main() {
  export PATH=$PATH:/jre/bin

  # We cannot use the IP of this node to performe the removal of this node of the cluster
  ZKHOST=$(sed -e "s/$HOST:3888//;s/^,//;s/,$//" < /opt/zookeeper/conf/server.list | cut -d ',' -f 1)
  ACTUAL_SERVERS=$(/opt/zookeeper/bin/zkCli.sh -server "$ZKHOST" config | grep "^server.")

  if echo "$ACTUAL_SERVERS" | grep -q "$HOST"; then
    echo "Removing $HOST server from zookeeper cluster"
    echo ""
    /opt/zookeeper/bin/zkCli.sh -server "$ZKHOST" reconfig -remove "$(cat /opt/zookeeper-data/data/myid)"
  fi
}
