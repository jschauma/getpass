# getpass - a Go module to get a password

[![Go
Reference](https://pkg.go.dev/badge/github.com/jschauma/getpass.svg)](https://pkg.go.dev/github.com/jschauma/getpass)

The `getpass` module provides a simple way to retrieve a password from
the user by specifying a number of different password sources.

A simple command-line tool using this function to
provide a generic helper command is also provided in
this directory.  Please see the [manual
page](https://github.com/jschauma/getpass/blob/main/doc/getpass.1.txt) for
details.

---

```
func Getpass(passfrom string) (pass string, err error)
```

Getpass retrieves a password from the user using a method defined by
the `passfrom` string.  The following methods are supported:

`cmd:command` -- Obtain the password by running the given command.  The
command will be passed to the shell for execution via `/bin/sh -c
'command'`.

`env:var` -- Obtain the password from the environment variable var.
Since the environment of other processes may be visible
via e.g. `ps(1)`, this option should be used with caution.

`fd:num` -- Obtain the password from the given file
descriptor.  Note that on many platforms this is
functionally equivalent to `file:/proc/$$/fd/num`.

`file:pathname` -- The first line of `pathname` is the password.
`pathname` need not refer to a regular file: it could for example refer
to a device or named pipe.  `pathname` undergoes standard "~" and
environment variable expansion.  Note that standard Unix file access
controls should be used to protect this file.

`keychain:name` -- Use the `security(1)` utility to retrieve the
password from the macOS keychain.

`lpass:name` -- Use the LastPass command-line client `lpass(1)` to
retrieve the named password.  You should previously have run `lpass
login` for this to work.

`op:name` -- Use the 1Password command-line client `op(1)` to retrieve
the named password.

`pass:password` -- The actual password is password.  Since the
password is visible to utilities such as `ps(1)` and possibly leaked
into the shell history file, this form should only be used where
security is not important.

`stdin` -- Read the password from stdin.  This is actually a convenience
alias for `fd:0`; on many platforms the same effect can be achieved via
`file:/dev/stdin`.

`tty[:prompt]` -- This is the default: `Getpass` will prompt the user on
the controlling tty using the provided `prompt`.  If no `prompt` is
provided, then `Getpass` will use "Password: ".

## Examples

```
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jschauma/getpass"
)

func main() {

	var pass string
	flag.StringVar(&pass, "p", "tty", "password method")
	flag.Parse()

	// Try out any of the options:
	p, err := getpass.Getpass(pass)
	if err != nil {
		log.Fatal("Unable to get password from user: ", err)
	} else {
		fmt.Printf("%s\n", p)
	}

	// Alternatively:

	// This will prompt the user to enter a password interactively,
	// using the default prompt.
	p, err = getpass.Getpass()
	if err != nil {
		log.Fatal("Unable to get password from user: ", err)
	} else {
		fmt.Printf("%s\n", p)
	}

	// Using a custom prompt:
	p, err = getpass.Getpass("tty:Please enter your secret passphrase: ")
	if err != nil {
		log.Fatal("Unable to get password from user: ", err)
	} else {
		fmt.Printf("%s\n", p)
	}

	// Using an environment variable:
	p, err = getpass.Getpass("env:MYSECRET")
	if err != nil {
		log.Fatal("Unable to get password from user: ", err)
	} else {
		fmt.Printf("%s\n", p)
	}

	// Using a file:
	p, err = getpass.Getpass("file:~/.secret")
	if err != nil {
		log.Fatal("Unable to get password from user: ", err)
	} else {
		fmt.Printf("%s\n", p)
	}

	// etc. etc.
}
```

---

See also:
* https://www.netmeister.org/blog/passing-passwords.html
* https://www.netmeister.org/blog/consistent-tools.html#passwords
