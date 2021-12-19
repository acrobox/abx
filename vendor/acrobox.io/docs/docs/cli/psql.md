# abx psql

Usage: `abx psql [OPTIONS] [ARGS...]`

Proxy options and arguments to the PostgreSQL container. See the PostgreSQL
documentation for `psql` for more information.

https://www.postgresql.org/docs/current/app-psql.html

## Examples

Query from the `example` database:

```sh
$ abx psql -c 'SELECT col FROM tablename' example
```

Enable the `citext` extension on the `example` database:

```sh
$ abx psql -c 'CREATE EXTENSION "citext"' example
```

The `psql` command is run as the `postgres` user internally. If you need to run
something as the user that owns the database created by `db/create` then use
the `psql` flag `-U acrobox`:

```sh
$ abx psql -U acrobox -c 'CREATE TABLE data ( id SERIAL PRIMARY KEY )' example
```

Export table `data` from the `example` database as a CSV document to `stdout`:

```sh
$ abx psql -c 'SELECT * FROM data' --csv example
```
