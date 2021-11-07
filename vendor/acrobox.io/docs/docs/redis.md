# Redis

Redis was selected as the in-memory key/value data structure store to be
tightly integrated with Acrobox. Redis is free, open source, well-documented,
and has a wide variety of community and commercial support options available.
Reach for Redis as a cache or message broker.

## Connection

Connect to the Redis server using the Data Store Name (DSN) connection string
available from the `ABX_CACHE_DSN` environment variable within containers.

The server is not publicly accessible. See `abx help redis-cli` for information
on how to connect to the Redis command line interface.

External connections can be made through Redis clients that support SSH key
authentication. See the output from `abx db/info` for connection information.

## Backups

Data is persisted to `/acrobox/redis/data.rdb` by default. It is recommended to
enable encrypted offsite backups. See `abx help backups` for more information.

## Upgrades

Acrobox will continually update Redis to the latest point release. Most Redis
application client libraries will automatically reconnect after the upgrade.
Major version upgrades must be performed manually.
