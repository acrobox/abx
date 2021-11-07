# Installation

## Source

```sh
$ go install acrobox.io/abx@latest
```

## Docker

Pull the image from the GitHub Container Registry.

```sh
$ docker pull ghcr.io/acrobox/abx
```

Create a `$ACROBOX_HOME` directory. This directory must be created before
running the image, otherwise Docker will create the volume mount as root.

```sh
$ mkdir -p "$XDG_CONFIG_HOME/abx"
```

Run the container as your host machine user.

```sh
$ docker run --rm -i -t -u $(id -u):$(id -g) -v "$XDG_CONFIG_HOME/abx":/machines ghcr.io/acrobox/abx help
```
