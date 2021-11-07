# abx destroy

Usage: `abx destroy [OPTIONS]`

Destroy an existing machine. **The machine and associated resources will be
destroyed and non-recoverable.**

Machines destroyed within 7 days of billing will be automatically refunded.

## Options

`-digitalocean-access-token` sets the DigitalOcean API access token. If a flag
is not provided, the environment variable `DIGITALOCEAN_ACCESS_TOKEN` will be
used by default. Your DigitalOcean access token will not be shared, stored, or
logged. Feel free to generate a new token specifically for Acrobox or update
your environment variables after use.

`-token` sets your acrobox.io token. Use of the environment variable
`ACROBOX_TOKEN` is recommended.

`-f` or `-force` to skip the confirmation prompt.
