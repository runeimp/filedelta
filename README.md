FileDelta
=========

1. Calculates the SHA-256 hash of a file
2. Can store the hash for later reference
3. Can check a file's current hash against the stored hash

This is a standalone tool for what many modern build systems already do, which is check the files hash to note if it has changed and needs to be rebuilt. This is to get around the old system which just checked the files modification time. Which isn't as realiable as one would hope.

```bash
$ filedelta source.go
583749f0278312c983b825d9a9025b40e8edad29c0ba368ce322a237d2098497 *source.go
$ filedelta store source.go
Hash stored
$ filedelta check source.go; echo "Exit Code: $?"
source.go: OK
Exit Code: 0
$ echo "Whaaaaat?" > source.go
$ filedelta check source.go; echo "Exit Code: $?"
source.go: ERROR
Exit Code: 1
```

Currently uses a local "cache" directory (in the current directory) for storing hash values. This will soon change.

This was inspired by a short discussion at https://github.com/casey/just/issues/424
