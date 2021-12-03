# Getting Started

Sign up at acrobox.io and generate an access token. If you've previously
generated one and didn't record it then navigate back to your account page
and regenerate the token. Set the `ACROBOX_TOKEN` environment variable on your
machine to this value.

```sh
$ export ACROBOX_TOKEN="hunter2"
```

Initialize a new Acrobox machine. This may take a few minutes, particularly
when provisioning the machine and associated resources and then again when
waiting for the initial setup to complete.

```sh
$ abx init
• Creating a new key pair.
• Payment authorization required.
  I authorize Acrobox to send instructions to the financial
  institution that issued my card to take payments from my card
  account in accordance with the terms of my agreement with you.
• Please type 'acrobox' to agree: acrobox
• Initializing with acrobox.io.
• Provisioning machine and associated resources.
• Writing machine configuration.
• Waiting for SSH connectivity.
• Waiting for machine setup.
• Waiting for service setup.
• Acrobox is ready.
```

See `abx` documentation for more details including important information about
retaining long term access to your machine.

At this point, you will require access to modify the DNS records for at least
one domain name. We're using `example.com` here. It is recommended to point
your DNS records to the Acrobox machine before adding the site.

If you use the `www` or naked domain, you need the following records, where
`IPv4` is replaced by the IPv4 address of the machine from `abx status`. The
Acrobox reverse proxy will take care of the redirects between them.

| Type | Name | Content                 |
| ---- | ---- | ----------------------- |
| A    | @    | *IPv4*                  |
| A    | www  | *IPv4*                  |
| CAA  | @    | 0 issue letsencrypt.org |

If you use a subdomain, you only need one A record:

| Type | Name | Content                 |
| ---- | ---- | ----------------------- |
| A    | sub  | *IPv4*                  |
| CAA  | @    | 0 issue letsencrypt.org |

These records may take a while to fully propagate across the internet. You can
do a quick check by running `dig +short example.com` replacing `example.com`
with your domain optionally including subdomain.

Configure a new container from the `hello-sigterm` image named and referred to
from other containers as `demo` and exposed at `example.com`. The hostname,
container name, and image name could have all been `example.com`. Different
names were chosen for clarity.

```sh
$ abx add -s example.com demo hello-sigterm
```

Build and deploy the `hello-sigterm` project or use one of your own.

```sh
$ git clone https://github.com/acrobox/hello-sigterm.git
$ cd hello-sigterm
$ docker build -t hello-sigterm .
$ abx deploy hello-sigterm
```

If everything went as expected, you should have a working app!

New applications are ready to connect to PostgreSQL and/or Redis.
See `abx help postgresql` and `abx help redis` for more information.

Many applications require background tasks to be run.
See `abx help add` and `abx help tasks` for more information.
