# abx remove

Usage: `abx remove [OPTIONS] NAME`

Remove an existing container configuration from the machine. If the container
is running, it will receive a `SIGKILL` signal. **Data will be destroyed.**

## Options

`-f` or `-force` to skip the confirmation prompt.
