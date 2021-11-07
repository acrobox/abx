# abx push

Usage: `abx push SOURCE... TARGET`

Copy files from the local machine to the host machine.

Only one source is allowed if the target exists as a file on the host machine.

Directories are copied recursively. Existing files will be overwritten.

## Examples

Use shell expansion to push all files and directories:

```sh
$ abx push ~/tmp/acrobox/* /acrobox
```

Explicitly push `foo` and `bar.txt` to `/acrobox/{foo,bar.txt}`:

```sh
$ abx push ~/tmp/dir1/foo ~/tmp/dir2/bar.txt /acrobox
```
