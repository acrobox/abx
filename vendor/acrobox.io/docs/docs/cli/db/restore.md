# abx db/restore

Usage: `abx db/restore [OPTIONS] NAME`

Restore a database from an archive file must be created by the automated
backups, `abx db/backup`, `pg_dump`, or a `pg_dump` compatible tool.

The database must already exist. Use `abx db/create` to create one if needed.
**All existing data will be lost.**

This command is functionally equivalent to
`pg_restore --clean --if-exists --no-acl --no-owner`. See the PostgreSQL
documentation for `pg_restore` for more information.

https://www.postgresql.org/docs/current/app-pgrestore.html

## Options

`-file` to specify the archive file. Defaults to `/data/postgres/NAME.dump`.

`-f` or `-force` to skip the confirmation prompt.
