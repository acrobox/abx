# abx cancel

Usage: `abx cancel [OPTIONS]`

Cancel an existing subscription. The machine and associated resources will be
preserved but Acrobox services will be paused.

Cancellations made within 7 days of billing will be automatically refunded but
service will be disabled immediately. Cancellations made thereafter will
maintain service until the end of the billing period.

## Options

`-token` sets your acrobox.io token. Use of the environment variable
`ACROBOX_TOKEN` is recommended.

`-f` or `-force` to skip the confirmation prompt.
