# abx run

Usage: `abx run IMAGE [COMMAND] [ARGS...]`

Run a one-off non-scheduled task.

These tasks do not need to have been configured by `abx add` but can still
communicate with those that have by container name. See `abx status` for a list
of configured containers and the names by which they can be referenced.

The host directory `/acrobox/IMAGE` will be mounted to `/data` on the container
with `acrobox` user permissions. For example, if `IMAGE` is `namespace/image`
then the host directory `/acrobox/namespace/image` will be mounted to `/data`
on the container.
