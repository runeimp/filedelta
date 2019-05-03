FileDelta v0.3.0
================

1. Calculates the SHA-256 hash of a file
2. Can store the hash for later reference
3. Can check a file's current hash against the stored hash

This is a standalone tool for what many modern build systems already do, which is check the files hash to note if it has changed and needs to be rebuilt. This is to get around the old system which just checked the files modification time. Which isn't as realiable as one would hope.

```bash
$ filedelta test/touched.txt
b9fa95a472cd1253bd7617700c44eb26b19d32fd32f9dd87a98976adf1c4fdd5 *test/touched.txt
$ filedelta store test/touched.txt
b9fa95a472cd1253bd7617700c44eb26b19d32fd32f9dd87a98976adf1c4fdd5 *test/touched.txt
$ filedelta check test/touched.txt; echo "Exit Code: $?"
test/touched.txt: OK
Exit Code: 0
$ echo "TOUCHED" > test/touched.txt
$ filedelta check test/touched.txt; echo "Exit Code: $?"
test/touched.txt: ERROR
Exit Code: 1
```

Hashes are now stored in `$HOME/.local/filedelta/cache/` in plain text files.

This was inspired by a short discussion at <https://github.com/casey/just/issues/424>

This tool uses [just][] instead of `make` for handling of task automation and builds. I highly recommend [just][] over `make`. It's awesome!  :smiley:


[just]: https://github.com/casey/just
