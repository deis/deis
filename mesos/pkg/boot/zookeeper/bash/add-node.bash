set -eo pipefail

# set debug based on envvar
[[ $DEBUG ]] && set -x

main() {
  export PATH=$PATH:/jre/bin

  cp /app/zoo.cfg /opt/zookeeper-data/zoo.cfg
  ln -s /opt/zookeeper-data/zoo.cfg /opt/zookeeper/conf/zoo.cfg

  cp /opt/zookeeper/conf/fleet-zoo_cfg.dynamic /opt/zookeeper-data/zoo_cfg.dynamic

  # # We need to add this node to the cluster if is not configured in the cluster
  # ZKHOST=$(sed -e "s/$HOST:3888//;s/^,//;s/,$//" < /opt/zookeeper/conf/server.list | cut -d ',' -f 1)

  # echo "adding $HOST as server to the zookeeper cluster"
  # echo ""
  # /opt/zookeeper/bin/zkCli.sh -server "$ZKHOST" reconfig -add "server.$(cat /opt/zookeeper-data/data/myid)=$HOST:2181:2888:participant;$HOST:3888"
}
