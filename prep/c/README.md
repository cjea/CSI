Loose clone of `ls`. More of a (~~bad~~ stripped down) DSL for listing file information.

## Build

```
$ make build
```

## Use

### Help

```
$ ./ls.out -h
```

Supports for columns of information: mode (`-m`), user (`-u`), group (`-g`), and size (`-s`). Mix and match them:

```
$ ./ls.out $(pwd)/
/Users/cjapel/misc/csi/prep/c/.
/Users/cjapel/misc/csi/prep/c/..
/Users/cjapel/misc/csi/prep/c/ls.c
/Users/cjapel/misc/csi/prep/c/Makefile
/Users/cjapel/misc/csi/prep/c/ls.out
/Users/cjapel/misc/csi/prep/c/README.md
/Users/cjapel/misc/csi/prep/c/pkg
```

```
$ ./ls.out -m $(pwd)/
drwxr-xr-x		/Users/cjapel/misc/csi/prep/c/.
drwxr-xr-x		/Users/cjapel/misc/csi/prep/c/..
-rw-r--r--		/Users/cjapel/misc/csi/prep/c/ls.c
-rw-r--r--		/Users/cjapel/misc/csi/prep/c/Makefile
-rwxr-xr-x		/Users/cjapel/misc/csi/prep/c/ls.out
-rw-r--r--		/Users/cjapel/misc/csi/prep/c/README.md
drwxr-xr-x		/Users/cjapel/misc/csi/prep/c/pkg
```

```
/ls.out -mugs $(pwd)/
drwxr-xr-x	cjapel	staff	288B		/Users/cjapel/misc/csi/prep/c/.
drwxr-xr-x	cjapel	staff	96B		/Users/cjapel/misc/csi/prep/c/..
-rw-r--r--	cjapel	staff	1692B		/Users/cjapel/misc/csi/prep/c/ls.c
-rw-r--r--	cjapel	staff	120B		/Users/cjapel/misc/csi/prep/c/Makefile
-rwxr-xr-x	cjapel	staff	19560B		/Users/cjapel/misc/csi/prep/c/ls.out
-rw-r--r--	cjapel	staff	40B		/Users/cjapel/misc/csi/prep/c/README.md
drwxr-xr-x	cjapel	staff	192B		/Users/cjapel/misc/csi/prep/c/pkg
```
