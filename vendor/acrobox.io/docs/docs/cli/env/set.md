# abx env/set

Usage: `abx env/set [OPTIONS] IMAGE KEY=value [KEY=value ...]`

Set one or more environment variables for the image.

## Options

`-e` or `-encoded` to read values as base64 encoded strings.

## Examples

```sh
$ abx env/set example.com EXTERNAL_SERVICE_TOKEN="hunter2"
```
