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
  I authorize Acrobox to charge my card in accordance with the terms of service.
• Please type 'acrobox' to agree: acrobox
• Initializing with acrobox.io.
• Provisioning machine and associated resources.
• Writing machine configuration.
• Waiting for SSH connectivity.
• Waiting for machine setup.
• Waiting for service setup.
• Acrobox is ready.
```

See `abx help` for more details including important information about
retaining long term access to your machine.

Before we go on, it is wise to enable backups. I like Backblaze B2 but you're
free to use other supported providers. Create a new bucket and generate
application keys. Set them as variables Acrobox will take care of the rest.

```sh
$ abx env/set acrobox/restic B2_ACCOUNT_ID="xxxxxxxxxxxxxxx0000000001"
$ abx env/set acrobox/restic B2_ACCOUNT_KEY="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
$ abx env/set acrobox/restic RESTIC_REPOSITORY="b2:bucketname"
$ abx env/set acrobox/restic RESTIC_PASSWORD="hunter2"
```

That's really it. Your repository will be automatically initialized on the
first scheduled run. See `abx help backups` for more details.

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

Configure a new container from the `hello-kitchen-sink` image named and
referred to from other containers as `demo` and exposed at `example.com`.
The hostname, container name, and image name could have all been `example.com`.
Different names were chosen for clarity.

```sh
$ abx add -s example.com demo hello-kitchen-sink
```

Build the `hello-kitchen-sink` project.

```sh
$ git clone https://github.com/acrobox/hello-kitchen-sink.git
$ cd hello-kitchen-sink
$ docker build -t hello-kitchen-sink .
```

We should create the database before the first deploy. This application doesn't
require a database on boot but many applications do.

```sh
$ abx db/create hello-kitchen-sink
$ abx psql -U acrobox -c 'CREATE TABLE data ( id SERIAL PRIMARY KEY )' hello-kitchen-sink
```

Now we're ready to deploy.

```sh
$ abx deploy hello-kitchen-sink
```

If everything went as expected, you should have a working app! This very basic
application prints the current value from a Redis key and the number of records
in a Postgres table, or plain text errors if encountered.

```sh
$ curl https://example.com
0
0
```

The response will always be the same. We can remedy this by adding a task to
increment a Redis key value and another task to add a record in a Postgres
table.

```sh
$ abx add -t "every 10 seconds" demo-incr hello-kitchen-sink incr
$ abx add -t "every 15 seconds" demo-touch hello-kitchen-sink touch
```

See `abx help add` and `abx help tasks` for more information.

Now if we wait some time, we can see that background tasks are running.

```sh
$ curl https://example.com
2
3
```

If the counts are getting too high for our liking, we can trigger a reset.

```sh
$ abx reload demo
```

Now that we have a general idea of the `abx` workflow, we can clean up
everything we created and move on to other projects.

```sh
$ abx remove demo
$ abx remove demo-incr
$ abx remove demo-touch
$ abx db/destroy hello-kitchen-sink
$ abx redis-cli del hello-kitchen-sink
$ abx ssh docker rmi hello-kitchen-sink
```
