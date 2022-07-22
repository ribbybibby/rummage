# rummage

Explore files in container images.

## Usage

List all the files in a container image:

```
$ rummage list alpine:latest
bin
bin/arch
bin/ash
bin/base64
bin/bbconfig
bin/busybox
bin/cat
bin/chattr
bin/chgrp
bin/chmod
bin/chown
bin/cp
bin/date
...
```

Use the `-l/--long` flag to print extended file information:

```
$ rummage list -l alpine:latest
drwxr-xr-x  0 0  0       2022-07-18 16:34:02 bin
Lrwxrwxrwx  0 0  0       2022-07-18 16:34:02 bin/arch -> /bin/busybox
Lrwxrwxrwx  0 0  0       2022-07-18 16:34:02 bin/ash -> /bin/busybox
Lrwxrwxrwx  0 0  0       2022-07-18 16:34:02 bin/base64 -> /bin/busybox
Lrwxrwxrwx  0 0  0       2022-07-18 16:34:02 bin/bbconfig -> /bin/busybox
-rwxr-xr-x  0 0  837272  2022-07-18 13:23:02 bin/busybox
...
```

By default, images are fetched directly from the remote registry. You can get
them from the docker daemon by setting the `-s/--source` flag:

```
$ rummage list -s daemon alpine:latest
```

Or a tar file:

```
$ crane pull alpine:latest alpine.tar
$ rummage list -s tarball alpine.tar
```

Layers pulled from a remote registry are cached to disk to speed up subsequent
lists. You can disable caching with the `--cache` flag:

```
$ rummage list --cache=false alpine:latest
```
