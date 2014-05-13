# logspout

A log router and HTTP interface for Docker container log streams, made to run inside Docker. Besides the routes you make, it's a stateless log appliance. It's not meant for managing log files or looking at history, just a means to get your logs out to live somewhere else, where they belong.

## Getting and running

Logspout is a (very small) Docker container, so you can just pull it from the index:

	$ docker pull progrium/logspout

When running logspout, it exposes port 8000 and needs two mounts. The first is the Docker Unix socket. The second is a directory to persist routes. We mount both with `-v`:

	$ docker run -d -P \
		-v=/var/run/docker.sock:/var/run/docker.sock \
		-v=/var/lib/logspout:/mnt/routes \
		progrium/logspout

Both need to be mounted in these specific paths inside the container, but where you keep the routes on the host could be anywhere. It could also be a regular Docker volume. If you don't mount a volume at `/mnt/routes`, it will only store routes in memory.

You can optionally pass an argument to install a catch-all route in the form `<type>://<addr>`. For example, to route all logs via syslog to `192.168.1.111:514`, run like this:

	$ docker run -d -P \
		-v=/var/run/docker.sock:/var/run/docker.sock \
		-v=/var/lib/logspout:/mnt/routes \
		progrium/logspout syslog://192.168.1.111:514

## HTTP API

### Streaming Endpoints

You can use these chunked transfer streaming endpoints for quick debugging with `curl` or for setting up easy TCP subscriptions to log sources. They also support WebSocket upgrades.

	GET /logs
	GET /logs/filter:<container-name-substring>
	GET /logs/id:<container-id>
	GET /logs/name:<container-name>

You can select specific log types from a source using a comma-delimited list in the query param `types`. Right now the only types are `stdout` and `stderr`, but when Docker properly takes over each container's syslog socket (or however they end up doing it), other types will be possible.

If you include a request `Accept: application/json` header, the output will be JSON objects including the name and ID of the container and the log type. Note that when upgrading to WebSocket, it will always use JSON.

Since `/logs` and `/logs/filter:<string>` endpoints can return logs from multiple source, they will by default return color-coded loglines prefixed with the name of the container. You can turn off the color escape codes with query param `colors=off` or the alternative is to stream the data in JSON format, which won't use colors or prefixes.


### Routing Resource

Routes let you configure logspout to hand-off logs to another system. Right now the only supported target type is via UDP `syslog`, but hey that's pretty much everything.

#### Creating a route

	POST /routes

Takes a JSON object like this:

	{
		"source": {
			"filter": "_db"
			"types": ["stdout"]
		},
		"target": {
			"type": "syslog",
			"addr": "logaggregator.service.consul"
			"append_tag": ".db"
		}
	}

The `source` field should be an object with `filter`, `name`, or `id` fields. You can specify specific log types with the `types` field to collect only `stdout` or `stderr`. If you don't specify `types`, it will route all types. 

To route all logs of all types on all containers, don't specify a `source`. 

The `append_tag` field of `target` is optional and specific to `syslog`. It lets you append to the tag of syslog packets for this route. By default the tag is `<container-name>`, so an `append_tag` value of `.app` would make the tag `<container-name>.app`.

And yes, you can just specify an IP and port for `addr`, but you can also specify a name that resolves via DNS to one or more SRV records. That means this works great with [Consul](http://www.consul.io/) for service discovery.

#### Listing routes

	GET /routes

Returns a JSON list of current routes:

	[
		{
			"id": "3631c027fb1b",
			"source": {
				"name": "mycontainer"
			},
			"target": {
				"type": "syslog",
				"addr": "192.168.1.111:514"
			}
		}
	]

#### Viewing a route

	GET /routes/<id>

Returns a JSON route object:

	{
		"id": "3631c027fb1b",
		"source": {
			"id": "a9efd0aeb470"
			"types": ["stderr"]
		},
		"target": {
			"type": "syslog",
			"addr": "192.168.1.111:514"
		}
	}

#### Deleting a route

	DELETE /routes/<id>

## Sponsor

This project was made possible by [DigitalOcean](http://digitalocean.com).

## License

BSD