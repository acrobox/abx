# Updates

## System

Alpine Linux updates run daily. The system will be rebooted if required.

## Docker

All containers shipped with Acrobox are updated, including itself. It is the
responsibility of the operator to redeploy application container updates.

It is often ideal to run updates during hours of low traffic. For many, that
window falls over night. Tasks are scheduled in UTC so 10am will fall between
2am and 6am across North America depending on daylight savings.

The update schedule can be augmented as required.

```sh
$ abx env/set acrobox/acroboxd UPDATE_SCHEDULE="10:00"
```

See `abx help tasks` for more information.

Updates will block until backups have completed.

Trigger a manual update with `abx update`.
