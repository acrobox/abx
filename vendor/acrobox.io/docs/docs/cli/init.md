# abx init

Usage: `abx init [OPTIONS]`

Initialize a new Acrobox machine. That's it! You're one command from having
provisioned secure, affordable, and low maintenance infrastructure --- complete
with tools and workflows to help you rapidly build, deploy, observe, and
iterate on multiple applications.

A new Droplet is created and provisioned with a Docker-based Alpine Linux
environment. Secure and lightweight by default, Acrobox is designed to squeeze
the most out of your virtual hardware.

Automatic updates for both the system and Docker images are applied daily.

An independently scalable block storage volume is created and mounted to
`/acrobox`. Acrobox configuration and site data is stored here. It is
recommended to enable encrypted offsite backups for this directory.

Only the following traffic will, by default, make it through a network layer
firewall to your server:

- `22/tcp` *(ssh)*
- `80/tcp` *(http to https redirects)*
- `443/tcp` *(https to serve applications)*

You retain full access to the infrastructure that is provisioned. The defaults
will cost you an additional $6.10 USD per month payable to DigitalOcean.

Host multiple sites on a single instance with the built in reverse proxy.
SSL/TLS certificates are issued by Let's Encrypt and yield an A+ on
the Qualys SSL Labs test.

Your Acrobox email address is registered with Let's Encrypt so that you may
receive notifications about your certificates. Use of Acrobox implies
acceptance of Let's Encrypt terms of service.

Automatic redirects from `http` to `https` and `www` to naked domain are
enabled by default to ensure your sites are accessible from less than perfect
manual URL entry.

View system resources, recent deploys, and monitor deployed sites and services
without committing additional time and resources until you're ready.

By initializing a new machine, you authorize Acrobox to charge your card in
accordance with the terms of service.

See `abx help legal` for details.

## Options

`-r` or `-region` sets the region for your server and block storage mount. The
default is `nyc1`.

Valid regions include:

- `nyc1` *(New York City, United States)*
- `ams3` *(Amsterdam, the Netherlands)*
- `sfo2` *(San Francisco, United States)*
- `sfo3` *(San Francisco, United States)*
- `sgp1` *(Singapore)*
- `lon1` *(London, United Kingdom)*
- `fra1` *(Frankfurt, Germany)*
- `tor1` *(Toronto, Canada)*
- `blr1` *(Bangalore, India)*

`-s` or `-size` sets the Droplet size. By default the cheapest DigitalOcean
Droplet size `s-1vcpu-1gb-intel` is provisioned. Use this option to provision
a larger instance. You can always scale up later, when the time is right.

`-d` or `-data-size` sets the block storage data volume size in gigabytes.
DigitalOcean supports block storage volumes from 1GB up to 16TB. The default is
1GB.

`-digitalocean-access-token` sets the DigitalOcean API access token. If a flag
is not provided, the environment variable `DIGITALOCEAN_ACCESS_TOKEN` will be
used by default. Your DigitalOcean access token will not be shared, stored, or
logged. Feel free to generate a new token specifically for Acrobox or update
your environment variables after use.

`-token` sets your acrobox.io token. Use of the environment variable
`ACROBOX_TOKEN` is recommended.

`-f` or `-force` to skip the confirmation prompt.
