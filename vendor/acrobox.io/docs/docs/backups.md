# Backups

restic was selected as the backup program to be tightly integrated with
Acrobox. restic is free, open source, well-documented, and performs
deduplicated, incremental, and encrypted cloud storage backups.

## Backup

The backup procedure runs as follows:

- All PostgreSQL databases are dumped to `/acrobox/postgres/DATABASE.dump`.
- `/acrobox/backup.sh` is executed if it exists.
- The entire `/acrobox` directory is backed up via restic if configured.

Directories marked with a special `CACHEDIR.TAG` file are excluded.
See https://bford.info/cachedir/ for more information.

Patterns specified in `/acrobox/restic-exclude.txt` will also be excluded.

https://restic.readthedocs.io/en/latest/040_backup.html#excluding-files

It is recommended to set backups to run automatically as a scheduled task. It
is often ideal to run backups during hours of low traffic. For many, that
window falls over night. By default, backups run at 10am UTC, somewhere between
2am and 6am across North America depending on daylight savings.

```sh
$ abx env/set acrobox/acroboxd BACKUP_SCHEDULE="10:00"
```

See `abx help tasks` for more information.

You need to set the following environment variables to enable restic support:

```sh
$ abx env/set acrobox/restic RESTIC_REPOSITORY="repo"
$ abx env/set acrobox/restic RESTIC_PASSWORD="hunter2"
```

Repositories hosted by external services may incur additional costs.

See the restic documentation for information about supported cloud storage
providers along with additional required or optional environment variables.

https://restic.readthedocs.io/en/latest/030_preparing_a_new_repo.html

See the pricing information for your cloud storage provider of choice and set
usage limits and/or notifications if appropriate.

Trigger a manual backup with `abx backup`.

## Restore

Install and use the `restic` command line application from any machine. It is
recommended to set the same restic environment variables locally to make the
command line program easier to use.

See the restic documentation for more information.

## Examples

View snapshots:

```sh
$ restic restore latest --target /tmp/acrobox
```

Download the latest snapshot of `/acrobox`:

```sh
$ restic restore latest --target /tmp/acrobox
```

Mount the latest snapshot of `/acrobox` to `/tmp/acrobox`:

```sh
$ restic mount /tmp/acrobox
```

Download the PostgreSQL database dump for the `example` database:

```sh
$ restic restore latest --target /tmp/acrobox \
  --path /acrobox/postgres/example.dump
```
