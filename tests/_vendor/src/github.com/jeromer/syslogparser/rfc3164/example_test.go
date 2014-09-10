package rfc3164_test

import (
	"fmt"
	"github.com/jeromer/syslogparser/rfc3164"
)

func ExampleNewParser() {
	b := "<34>Oct 11 22:14:15 mymachine su: 'su root' failed for lonvick on /dev/pts/8"
	buff := []byte(b)

	p := rfc3164.NewParser(buff)
	err := p.Parse()
	if err != nil {
		panic(err)
	}

	fmt.Println(p.Dump())
}
