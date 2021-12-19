# abx add

Usage: `abx add [OPTIONS] NAME IMAGE [COMMAND] [ARGS...]`

Configure a new container on the machine.

Container configuration is integral to `abx` functionality. Most `abx` commands
only affect containers configured with the `abx add` command. All other
containers are to be manually managed by the operator.

Configuration is stored within `/acrobox/config.json`. Manual modification of
this file is not recommended.

The provided `NAME` is normalized to a valid container name beginning with an
alphanumeric character and followed by one or more alphanumeric characters,
underscores, periods, or hyphens. The name must be unique.

Containers can resolve to each other by the chosen `NAME`. It may be convenient
to name a container after the domain name they publicly resolve to but it may
raise issues when connecting from other containers.

See `abx status` for a list of configured containers.

The configuration applies to the most recently deployed `IMAGE` so that the
service, site, or task only needs to be configured once. Image tags are
ignored. See `abx help deploy` for more information.

At most one of `-t` and `-s` may be specified.

The `-t` option specifies the temporal expression on which to run the container
as a scheduled task. See `abx help tasks` for more information. All other
containers are configured to run as a persistent background service.

Use `-s` to register the service with the reverse proxy such that incoming
requests to the given host and optional path are routed to the port within the
container. The proxy will take care of HTTP/HTTPS and naked/www redirects. See
the proxy documentation for more details.

https://pkg.go.dev/github.com/acrobox/proxy#Manager.Add

## Options

`-t` or `-task` to schedule the container as a task.

`-s` or `-site` to specify the host and path to expose the container through.

`-p` or `-port` to specify the site container port. Defaults to "8080".
