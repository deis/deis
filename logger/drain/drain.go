package drain

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
)

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
