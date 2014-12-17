Using this library you can easy implement your own syslog server that:

1. Can listen on multiple UDP ports and unix domain sockets.

2. Can pass parsed syslog messages to your own handlers so your code can analyze
and respond for them.

See [documentation](http://gopkgdoc.appspot.com/pkg/github.com/ziutek/syslog)
and [example server](https://github.com/ziutek/syslog/blob/master/example_server/main.go).
