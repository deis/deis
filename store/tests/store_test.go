package tests

import (
  "fmt"
  "testing"
  "time"

  "github.com/deis/deis/tests/dockercli"
  "github.com/deis/deis/tests/etcdutils"
  "github.com/deis/deis/tests/utils"
)

func TestStore(t *testing.T) {
  var err error
  setkeys := []string{
  }
  setdir := []string{
    "/deis/store",
  }
  // Set up etcd, which will be used by all containers
  tag, etcdPort := utils.BuildTag(), utils.RandomPort()
  etcdName := "deis-etcd-" + tag
  cli, stdout, stdoutPipe := dockercli.NewClient()
  dockercli.RunTestEtcd(t, etcdName, etcdPort)
  defer cli.CmdRm("-f", etcdName)
  host := utils.HostAddress()

  // test deis-store-monitor
  handler := etcdutils.InitEtcd(setdir, setkeys, etcdPort)
  etcdutils.PublishEtcd(t, handler)
  etcdutils.SetEtcd(t,
    []string{"/deis/store/hosts/"+host},
    []string{"boot2docker"},
    handler.C)

  fmt.Printf("--- Run deis/store-monitor:%s at %s:%s\n", tag, host, "6789")
  name := "deis-store-monitor-" + tag
  defer cli.CmdRm("-f", name)
  go func() {
    _ = cli.CmdRm("-f", name)
    err = dockercli.RunContainer(cli,
      "--name", name,
      "--rm",
      "-p", "6789",
      "-e", "HOST="+host,
      "-e", "ETCD_PORT="+etcdPort,
      "-e", "NUM_STORES=1",
      "--net=host",
      "deis/store-monitor:"+tag)
  }()
  dockercli.PrintToStdout(t, stdout, stdoutPipe, "monmap e1: 1 mons at")
  if err != nil {
    t.Fatal(err)
  }
  // FIXME: Wait until etcd keys are published
  time.Sleep(5000 * time.Millisecond)
  dockercli.DeisServiceTest(t, name, "6789", "tcp")
  etcdutils.VerifyEtcdKey(t, "/deis/store/conf", etcdPort)
  etcdutils.VerifyEtcdKey(t, "/deis/store/monKeyring", etcdPort)
  etcdutils.VerifyEtcdKey(t, "/deis/store/adminKeyring", etcdPort)
  etcdutils.VerifyEtcdValue(t, "/deis/store/monSetupComplete", "youBetcha", etcdPort)

  // test deis-store-daemon
  fmt.Printf("--- Run deis/store-daemon:%s at %s:%s\n", tag, host, "6800")
  name = "deis-store-daemon-" + tag
  defer cli.CmdRm("-f", name)
  go func() {
    _ = cli.CmdRm("-f", name)
    err = dockercli.RunContainer(cli,
      "--name", name,
      "--rm",
      "-p", "6800",
      "-e", "HOST="+host,
      "-e", "ETCD_PORT="+etcdPort,
      "--net=host",
      "deis/store-daemon:"+tag)
  }()
  dockercli.PrintToStdout(t, stdout, stdoutPipe, "Booting")
  if err != nil {
    t.Fatal(err)
  }
  // FIXME: Wait until etcd keys are published
  time.Sleep(5000 * time.Millisecond)
  dockercli.DeisServiceTest(t, name, "6800", "tcp")
  etcdutils.VerifyEtcdValue(t, "/deis/store/osds/"+host, "0", etcdPort)

  // test deis-gateway
  port := utils.RandomPort()
  fmt.Printf("--- Run deis/store-gateway:%s at %s:%s\n", tag, host, port)
  name = "deis-store-gateway-" + tag
  defer cli.CmdRm("-f", name)
  go func() {
    _ = cli.CmdRm("-f", name)
    err = dockercli.RunContainer(cli,
      "--name", name,
      "--rm",
      "-p", port+":8888",
      "-e", "HOST="+host,
      "-e", "EXTERNAL_PORT="+port,
      "-e", "ETCD_PORT="+etcdPort,
      "deis/store-gateway:"+tag)
  }()
  dockercli.PrintToStdout(t, stdout, stdoutPipe, "Booting")
  if err != nil {
    t.Fatal(err)
  }
  // FIXME: Wait until etcd keys are published
  time.Sleep(5000 * time.Millisecond)
  dockercli.DeisServiceTest(t, name, port, "http")
  etcdutils.VerifyEtcdValue(t, "/deis/store/gateway/host", host, etcdPort)
  etcdutils.VerifyEtcdValue(t, "/deis/store/gateway/port", port, etcdPort)
  etcdutils.VerifyEtcdKey(t, "/deis/store/gatewayKeyring", etcdPort)
  etcdutils.VerifyEtcdKey(t, "/deis/store/gateway/accessKey", etcdPort)
  etcdutils.VerifyEtcdKey(t, "/deis/store/gateway/secretKey", etcdPort)
}
