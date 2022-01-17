# abx start

Usage: `abx start NAME`

Start a stopped container.

The container `NAME` must have been configured through `abx add`. Container
names `postgres` and `redis` are also supported.

Services and sites are run as daemons and will be automatically restarted
unless manually stopped. Tasks are run as one-off tasks and will not run if the
previous invocation is still running.

Containers managed by `abx` are connected to an isolated bridge network named
`acrobox`, allowing containers to resolve each other by name. See `abx status`
for a list of configured containers and the names by which they can be
referenced.

The host directory `/acrobox/IMAGE` will be mounted to `/data` on the container
with `acrobox` user permissions, where `IMAGE` is the image name previously
configured by `abx add`.

For example, if `IMAGE` is `namespace/image` then the host directory
`/acrobox/namespace/image` will be mounted to `/data` on the container.

Note that the directory is named using the underlying image name, not the
container name. This way, a single directory may be shared across multiple
services, sites, and tasks.

Likewise, environment variables set by `abx env/set` are bound to the image
and applied to all Acrobox containers that run the image.

Similarly, the `ABX_STORE_DSN` environment variable is set to a PostgreSQL
connection string using the underlying image name. This database is not created
automatically. Use `abx db/create` to create it.  See `abx help postgresql` for
more information.

The following environment variables are injected into the container runtime:

| Name            | Description                  |
| --------------- | ---------------------------- |
| `ABX_IMAGE`     | container image              |
| `ABX_NAME`      | container name               |
| `ABX_SITE`      | container site               |
| `ABX_PORT`      | container port               |
| `ABX_SCHEDULE`  | container schedule           |
| `ABX_STORE_DSN` | PostgreSQL connection string |
| `ABX_CACHE_DSN` | Redis connection string      |
