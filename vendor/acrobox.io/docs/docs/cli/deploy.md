# abx deploy

Usage: `abx deploy [OPTIONS] IMAGE`

Deploy a Docker image from your local machine to the host machine.

If an image tag is provided, it is ignored. The `latest` tag, or the last build
or tag that ran without an explicit tag, is deployed instead.

If there exists one or more containers previously configured by `abx add` that
are running services or sites using the given image then they will be
restarted. If this is the first deploy, they will simply be started. See `abx
help start` for details on how containers are run.

Restarts send a `SIGTERM` signal to the container. It is the responsibility of
the application to handle the signal. The intended use is to trigger a graceful
stop. If the application doesn't stop before the grace period, a `SIGKILL`
signal is sent to force kill the container. In any case, the container is
restarted.

## Options

`-f` or `-force` to immediately send `SIGKILL`.

`-t` or `-time` to set the grace period between sending `SIGTERM` and
`SIGKILL` on the container. Defaults to 10. Ignored if force is enabled.
