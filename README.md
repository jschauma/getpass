# getpass - a Go module to get a password

The `getpass` module provides a simple way to retrieve
a password from the user by specifying a number of
different password sources:

`env:var` -- Obtain the password from the environment variable var.
Since the environment of other processes may be visible
via e.g. ps(1), this option should be used with caution.

`file:pathname` -- The first line of pathname is the password.  pathname need
not refer to a regular file: it could for example refer to
a device or named pipe.  Note that standard Unix file
access controls should be used to protect this file.

`keychain:name` -- Use the security(1) utility to retrieve the
password from the macOS keychain.

`lpass:name` -- Use the LastPass command-line client lpass(1) to
retrieve the named password.  You should previously have
run 'lpass login' for this to work.

`op:name` -- Use the 1Password command-line client op(1) to
retrieve the named password.

`pass:password` -- The actual password is password.  Since the password is
visible to utilities such as ps(1) and possibly leaked
into the shell history file, this form should only be
used where security is not important.

If no password retrieval method is specified, then
`getpass.Getpass` will prompt the user on the
controlling tty using the provided prompt.

See also:
* https://www.netmeister.org/blog/passing-passwords.html
