# abx reload

Usage: `abx reload [OPTIONS] NAME`

Send a SIGUSR1 signal to the container. It is the responsibility of the
application to respond to the signal. The intended use is to reload
configuration without a hard restart.
