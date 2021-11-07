# Acrobox CLI

`abx` is the Acrobox CLI program.

All options can be set as environment variables prefixed with `ACROBOX_` unless
otherwise specified.

Application configuration can be found in `$ACROBOX_HOME`. You are responsible
for backing up this data to prevent loss of access. If `$ACROBOX_HOME` is not
set, the first applicable of following locations will be used:

- `$XDG_CONFIG_HOME/acrobox` on Unix systems
- `$HOME/.config/acrobox` on Unix systems when `$XDG_CONFIG_HOME` is empty
- `$HOME/Library/Application Support/acrobox` on macOS
- `%AppData%/acrobox` on Windows
- `$home/lib/acrobox` on Plan 9

Errors yield a non-zero exit status.

## Options

`-h` or `-host` specifies the machine hostname. Defaults to "acrobox".

`-v` or `-verbose` enables verbose output.

## Commands

### Machine

`init` provisions a new machine.

`cancel` cancels an existing subscription.

`renew` renews an existing subscription.

`destroy` destroys a machine.

`ssh` logs into the machine.

`push` copies files from the local machine to the host machine.

`pull` copies files from the host machine to the local machine.

`status` displays machine status information.

`db/info` prints database connection information.

`db/list` displays configured databases.

`db/create` creates a new database.

`db/destroy` destroys an existing database.

`psql` proxies to psql on the machine.

`redis-cli` proxies to redis-cli on the machine.

`backup` triggers a manual backup.

`update` triggers a manual update.

### Containers

`add` configures a new container on the machine.

`remove` removes an existing container configuration.

`deploy` deploys an image.

`list` displays configured containers.

`show` displays container information.

`logs` displays container logs.

`stop` stops a container.

`start` starts a stopped container.

`reload` signals the container to reload configuration.

`restart` restarts a container.

`env/set` sets one or more image environment variables.

`env/get` prints an image environment variable.

`env/all` prints all image environment variables.

`env/del` deletes one or more image environment variables.

### Program

`help` displays information on commands and additional topics.

`version` prints the application version.

## Additional Help Topics

`installation` displays the installation instructions.

`getting-started` displays information on provisioning and deployment.

`security` displays information about the security considerations.

`firewall` displays information on the network layer firewall setup.

`backups` displays information on configuring and maintaining backups.

`updates` displays information on updates.

`postgresql` displays information on managing the PostgreSQL database.

`redis` displays information on managing the Redis cache server.

`tasks` displays information on managing scheduled tasks.

`legal` displays the privacy policy and terms of service.
