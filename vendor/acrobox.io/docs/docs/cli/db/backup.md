# abx db/backup

Usage: `abx db/backup [OPTIONS] NAME`

Extract a database to a PostgreSQL archive file.

This command is functionally equivalent to `pg_dump --format=custom`. See the
PostgreSQL documentation for `pg_dump` for more information.

https://www.postgresql.org/docs/current/app-pgdump.html

## Options

`-file` to specify the archive file. Defaults to `/data/postgres/backups/NAME.dump`.
