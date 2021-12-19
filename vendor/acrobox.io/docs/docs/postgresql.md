# PostgreSQL

PostgreSQL was selected as the data store to be tightly integrated with
Acrobox. PostgreSQL is free, open source, well-documented, and has a wide
variety of community and commercial support options available.

## Manage

Databases can be created and destroyed through the Acrobox CLI.

See `abx help db/create` and `abx help db/destroy` for information on how to
create and destroy databases.

## Connection

Connect to the PostgreSQL server using the Data Store Name (DSN) connection
string available from the `ABX_STORE_DSN` environment variable within
containers. The environment variable is always present but the database is not
created automatically. Use `abx db/create` to create it.

See the PostgreSQL documentation for more information about connection strings.

https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING

The database is not publicly accessible. See `abx help psql` for information on
how to connect to the PostgreSQL interactive terminal or to run ad hoc scripts
or commands from the command line.

External connections can be made through a PostgreSQL client that supports SSH
key authentication. See the output from `abx db/info` for connection details.

## Configuration

### Server

Alter server configuration at `/acrobox/postgres/postgresql.conf` as required.

See the PostgreSQL documentation for more information.

https://www.postgresql.org/docs/14/runtime-config.html

### psql

Alter `psql` configuration at `/acrobox/postgres/psqlrc` as required.

See the PostgreSQL documentation for more information.

https://www.postgresql.org/docs/current/app-psql.html

## Backups

Daily database dumps to `/acrobox/postgres/backups/NAME.dump` are enabled by
default. It is recommended to enable encrypted offsite backups. See `abx help
backups` for more information.

## Upgrades

Acrobox will continually update PostgreSQL to the latest point release. Most
PostgreSQL application client libraries will automatically reconnect after the
upgrade. Major version upgrades must be performed manually.
