package syslog

import (
    "reflect"
    "testing"
)

func expect(t *testing.T, a interface{}, b interface{}) {
    if a != b {
        t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
    }
}

func TestParsePacket(t *testing.T) {
    buf := []byte("2014-08-31 07:38:55 uuid[web.3]: 2014-08-31 07:38:47 [20] [INFO] Starting gunicorn 18.0")

    // test will fail if a panic is raised
    m := ParsePacket(buf)
    // ensure that if the log line is missing a priority, everything will be all right
    expect(t, int(m.Priority), 0)
    expect(t, m.Content, "2014-08-31 07:38:47 [20] [INFO] Starting gunicorn 18.0")
    expect(t, m.Hostname, "")
    expect(t, m.Tag, "uuid[web.3]")
}
