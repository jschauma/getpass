/*
 * Originally written by Jan Schaumann
 * <jschauma@netmeister.org> in May 2022.
 *
 */

// Package getpass provides a simple way to retrieve a
// password.
package getpass

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
)

// Getpass retrieves a password from the user using a
// method defined by the 'passfrom' string.  The
// following methods are supported:
//
//  env:var        Obtain the password from the environment variable var.
//                 Since the environment of other processes may be visible
//                 via e.g. ps(1), this option should be used with caution.
//
//  file:pathname  The first line of pathname is the password.  pathname need
//                 not refer to a regular file: it could for example refer to
//                 a device or named pipe.  Note that standard Unix file
//                 access controls should be used to protect this file.
//
//  keychain:name  Use the security(1) utility to retrieve the
//                 password from the macOS keychain.
//
//  lpass:name     Use the LastPass command-line client lpass(1) to
//                 retrieve the named password.  You should previously have
//                 run 'lpass login' for this to work.
//
//  op:name        Use the 1Password command-line client op(1) to
//                 retrieve the named password.
//
//  pass:password  The actual password is password.  Since the password is
//                 visible to utilities such as ps(1) and possibly leaked
//                 into the shell history file, this form should only be
//                 used where security is not important.
//
// If no password retrieval method is specified, then
// Getpass will prompt the user on the controlling tty
// using the provided prompt.
func Getpass(passfrom, prompt string) (pass string, err error) {
	var source string
	var passin []string
	errMsg := "invalid password source"
	if len(passfrom) == 0 {
		source = "tty"
	} else {
		passin = strings.SplitN(passfrom, ":", 2)
		if len(passin) < 2 {
			return "", errors.New(errMsg)
		}
		source = passin[0]
	}

	switch source {
	case "env":
		return getpassFromEnv(passin[1])
	case "file":
		return getpassFromFile(passin[1])
	case "keychain":
		return getpassFromKeychain(passin[1])
	case "lastpass":
		fallthrough
	case "lpass":
		return getpassFromLastpass(passin[1])
	case "onepass":
		fallthrough
	case "op":
		return getpassFromOnepass(passin[1])
	case "pass":
		return passin[1], nil
	case "tty":
		return getpassFromUser(prompt)
	default:
		return "", errors.New(errMsg)
	}

	return pass, nil
}

func getpassFromEnv(varname string) (pass string, err error) {
	errMsg := fmt.Sprintf("environment variable '%v' not set", varname)
	pass = os.Getenv(varname)
	if len(pass) < 1 {
		return "", errors.New(errMsg)
	}
	return pass, nil
}

func getpassFromFile(fname string) (pass string, err error) {
	file, err := os.Open(fname)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to open '%s': %v", fname, err)
		return "", errors.New(errMsg)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pass = scanner.Text()
	}

	return pass, nil
}

func getpassFromKeychain(entry string) (pass string, err error) {
	cmd := []string{"security", "find-generic-password", "-s", entry, "-w"}
	out, err := runCommand(cmd, "", false)
	if err != nil {
		return "", err
	}
	return out, nil
}

func getpassFromLastpass(entry string) (pass string, err error) {
	cmd := []string{"lpass", "show", entry, "--password"}
	out, err := runCommand(cmd, "", false)
	if err != nil {
		return "", err
	}
	return out, nil
}

func getpassFromOnepass(entry string) (pass string, err error) {
	cmd := []string{"op", "item", "get", entry, "--fields", "password"}
	out, err := runCommand(cmd, "", false)
	if err != nil {
		return "", err
	}
	return out, nil
}

func getpassFromUser(prompt string) (pass string, err error) {
	dev_tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return "", err
	}

	fmt.Fprintf(dev_tty, prompt)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		err = stty("echo")
		if err != nil {
			os.Exit(1)
		}
	}()

	err = stty("-echo")
	if err != nil {
		return "", err
	}

	input := bufio.NewReader(dev_tty)
	b, err := input.ReadBytes('\n')
	if err != nil {
		return "", err
	}
	pass = string(b)

	err = stty("echo")
	if err != nil {
		return "", err
	}
	fmt.Fprintf(dev_tty, "\n")

	return string(pass[:len(pass)-1]), nil
}

func runCommand(args []string, stdinData string, need_tty bool) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if need_tty {
		dev_tty, err := os.Open("/dev/tty")
		if err != nil {
			return "", err
		}
		cmd.Stdin = dev_tty
	} else if len(stdinData) > 0 {
		var b bytes.Buffer
		b.Write([]byte(stdinData))
		cmd.Stdin = &b
	}

	err := cmd.Run()
	if err != nil {
		errMsg := fmt.Sprintf("unable to run '%v':\n%v\n%v", strings.Join(args, " "), stderr.String(), err)
		return "", errors.New(errMsg)

	}
	return strings.TrimSpace(stdout.String()), nil
}

func stty(arg string) (err error) {
	flag := "-f"
	if runtime.GOOS == "linux" {
		flag = "-F"
	}

	err = exec.Command("/bin/stty", flag, "/dev/tty", arg).Run()
	if err != nil {
		errMsg := fmt.Sprintf("Unable to run stty on /dev/tty: %v", err)
		return errors.New(errMsg)
	}

	return nil
}
