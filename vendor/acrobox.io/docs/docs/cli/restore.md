# abx restore

Usage: `abx restore [OPTIONS] [ARGS...]`

Restore the contents from the configured restic repository.

This is potentially a destructive action as existing files may be overwritten.

If no arguments are passed, then `latest --target /data` is used to restore
the entire `/acrobox` block storage volume.

This command is restricted to targets under `/data`.

See `abx help backups` for more information on backups.

## Options

`-f` or `-force` to skip the confirmation prompt.
