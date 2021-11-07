# Security

Alpine Linux was selected as the base operating system on which Acrobox has
been built. Alpine Linux is a relatively small distribution designed to be
secure, simple, and resource efficient.

## Infrastructure

Consider adding a strong password and enabling 2FA on third party services
including but not limited to your domain registrar, DNS provider, and
DigitalOcean. Ideally, all three of these should be separate entities.

## Connectivity

Connectivity is made through a hardened SSH server. A new RSA key pair is
created for each machine and the SSH server has been configured to only allow
non-root public key connections. Strict host key verification is enforced for
all connections.

## Disclosure

I take security very seriously, but no system is perfect. Please report flaws
directly to hello@acrobox.io rather than publicly online. You will receive
credit, should you desire, for your research when the vulnerability is fixed.
