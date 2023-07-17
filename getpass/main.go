// Originally written by Jan Schaumann
// <jschauma@netmeister.org> in July 2023.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jschauma/getpass"
)

const PROGNAME = "getpass"
const VERSION = "0.2.2"

func usage() {
	help := `Usage: getpass [-Vh] [passin]
        -V  print version numbe and exit
        -h  print this help and exit

passin may be one of:
cmd:command, env:var, fd:num, file:pathname, keychain:name,
lpass:name, op:name, pass:password, stdin, tty[:prompt]
`
	fmt.Printf(help)
}

func main() {

	printVersion := false
	flag.BoolVar(&printVersion, "V", false, "print help")

	printHelp := false
	flag.BoolVar(&printHelp, "h", false, "print help")

	flag.Parse()
	if printHelp {
		usage()
		os.Exit(0)
	}

	if printVersion {
		fmt.Printf("%s version %s\n", PROGNAME, VERSION)
		os.Exit(0)
	}

	if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-Vh] [passin]\n", PROGNAME)
		os.Exit(1)
	}

	pass := "tty"
	if len(os.Args) == 2 {
		pass = os.Args[1]
	}

	p, err := getpass.Getpass(pass)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: Unable to getpass using '%s: %s\n",
			PROGNAME, pass, err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", p)
}
