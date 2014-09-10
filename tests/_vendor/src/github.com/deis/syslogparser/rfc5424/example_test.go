package rfc5424_test

import (
	"fmt"
	"github.com/jeromer/syslogparser/rfc5424"
)

func ExampleNewParser() {
	b := `<165>1 2003-10-11T22:14:15.003Z mymachine.example.com evntslog - ID47 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] An application event log entry...`
	buff := []byte(b)

	p := rfc5424.NewParser(buff)
	err := p.Parse()
	if err != nil {
		panic(err)
	}

	fmt.Println(p.Dump())
}
