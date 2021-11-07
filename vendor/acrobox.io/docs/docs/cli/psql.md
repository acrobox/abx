# abx psql

Usage: `abx psql [OPTIONS] [ARGS...]`

Proxy options and arguments to the PostgreSQL container. See the PostgreSQL
documentation for `psql` for more information.

https://www.postgresql.org/docs/current/app-psql.html

## Examples

Query from the `example` database:

```sh
$ abx psql -c 'SELECT col FROM tablename;' example
```

Enable the `citext` extension on the `example` database:

```sh
$ abx psql -U postgres -c 'CREATE EXTENSION "citext"' example
```
