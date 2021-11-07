# abx restart

Usage: `abx restart [OPTIONS] NAME`

Send a `SIGTERM` signal to the container. It is the responsibility of the
application to handle the signal. The intended use is to trigger a graceful
stop. If the application doesn't stop before the grace period, a `SIGKILL`
signal is sent to force kill the container. In any case, the container is
restarted.

See `abx help start` for details on how containers are run.

## Options

`-f` or `-force` to immediately send `SIGKILL`.

`-t` or `-time` to set the grace period between sending `SIGTERM` and
`SIGKILL` on the container. Defaults to 10. Ignored if force is enabled.
