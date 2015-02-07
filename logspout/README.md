# logspout

A log router for Docker container output that runs entirely inside Docker. It attaches to all containers on a host, then routes their logs wherever you want.

It's a 100% stateless log appliance (unless you persist routes). It's not meant for managing log files or looking at history. It is just a means to get your logs out to live somewhere else, where they belong.

For now it only captures stdout and stderr, but soon Docker will let us hook into more ... perhaps getting everything from every container's /dev/log.

## Getting logspout

Logspout is a very small Docker container, so you can just pull it from the index:

	$ docker pull deis/logspout

## Using logspout

#### Route all container output to remote syslog

The simplest way to use logspout is to just take all logs and ship to a remote syslog. Just pass a default syslog target URI as the command. Also, we always mount the Docker Unix socket with `-v` to `/tmp/docker.sock`:

	$ docker run -v=/var/run/docker.sock:/tmp/docker.sock deis/logspout /bin/logspout syslog://logs.papertrailapp.com:55555

If deis/logspout is deployed on Deis, it will connect automatically to deis-logger via service discovery.

#### Inspect log streams using curl

Whether or not you run it with a default routing target, if you publish its port 8000, you can connect with curl to see your local aggregated logs in realtime.

	$ docker run -d -p 8000:8000 \
		-v=/var/run/docker.sock:/tmp/docker.sock \
		deis/logspout
	$ curl $(docker port `docker ps -lq` 8000)/logs

You should see a nicely colored stream of all your container logs. You can filter by container name, log type, and more. You can also get JSON objects, or you can upgrade to WebSocket and get JSON logs in your browser.

See [Streaming Endpoints](#streaming-endpoints) for all options.

#### Create custom routes via HTTP

Along with streaming endpoints, logspout also exposes a `/routes` resource to create and manage routes.

	$ curl $(docker port `docker ps -lq` 8000)/logs -X POST \
		-d '{"source": {"filter": "db", "types": ["stderr"]}, target": {"type": "syslog", "addr": "logs.papertrailapp.com:55555"}}'

That example creates a new syslog route to [Papertrail](https://papertrailapp.com) of only `stderr` for containers with `db` in their name.

By default, routes are ephemeral. But if you mount a volume to `/mnt/routes`, they will be persisted to disk.

See [Routes Resource](#routes-resource) for all options.

#### Using a custom timestamp format

By default, logspout will use the timestamp format `2006-01-02T15:04:05MST`. A custom format can be specified by setting the `DATETIME_FORMAT` environment variable.

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


### Routes Resource

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
