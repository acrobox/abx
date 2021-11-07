# Firewall

Only the following traffic will, by default, make it through a network layer
firewall to your server:

- `22/tcp` *(ssh)*
- `80/tcp` *(http to https redirects)*
- `443/tcp` *(https to serve applications)*

You retain full access to the infrastructure that is provisioned. Feel free to
augment these firewall rules from the DigitalOcean dashboard if needed.
