package drain

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"

	"github.com/coreos/go-etcd/etcd"
)

func GetDrain() string {
	host := getopt("HOST", "127.0.0.1")

	etcdPort := getopt("ETCD_PORT", "4001")
	etcdPath := getopt("ETCD_PATH", "/deis/logs")

	client := etcd.NewClient([]string{"http://" + host + ":" + etcdPort})

	s, err := client.Get(etcdPath+"/drain", true, false)
	if err != nil {
		return ""
	}

	return s.Node.Value
}

func SendToDrain(m string, drain string) error {
	u, err := url.Parse(drain)
	if err != nil {
		log.Fatal(err)
	}
	uri := u.Host + u.Path
	switch u.Scheme {
	case "syslog":
		sendToSyslogDrain(m, uri)
	default:
		log.Println(u.Scheme + " drain type is not implemented.")
	}
	return nil
}

func sendToSyslogDrain(m string, drain string) error {
	conn, err := net.Dial("udp", drain)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	fmt.Fprintf(conn, m)
	return nil
}

func getopt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		value = dfault
	}
	return value
}
