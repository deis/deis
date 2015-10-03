package drain

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/url"
	"net/http"
	"os"
	"strings"
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
	case "https":
		sendToHttpsDrain(m, uri)
	default:
		log.Println(u.Scheme + " drain type is not implemented.")
	}
	return nil
}

func sendToHttpsDrain(m string, drain string) error {
	buf := strings.NewReader(m)

	tr := &http.Transport{
  	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
  }

  client := &http.Client{Transport: tr}

	// Skip logging error here since any log line we write will be cycled
	// infinitely through the drain if there is an ongoing error
  resp, err := client.Post(("https://" + drain), "text/plain", buf)
	resp.Body.Close()
	return err
}

func sendToSyslogDrain(m string, drain string) error {
	conn, err := net.Dial("udp", drain)
	if err != nil {
		log.Print(err)
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
