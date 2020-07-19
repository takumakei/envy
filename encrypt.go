package main

import (
	"flag"
	"fmt"
	"os"
)

var encryptUsage = `usage: envy encrypt [options]

  encrypt the file

options:
  -h           show this message
  -f <file>    input filename to encrypt
               (default "-" stdin)
  -p <env>     enviromnent variable name of the password
               (default "ENVY_PASSWORD")
  -o <file>    output filename
               (default "-" stdout)
  -Y           allow to overwrite the output file
`

func encrypt(args []string) int {
	var out string
	var file string
	var passenv string
	var overwrite bool
	commandLine := flag.NewFlagSet("encrypt", flag.ExitOnError)
	commandLine.StringVar(&out, "o", out, "output file")
	commandLine.StringVar(&file, "f", file, "input file")
	commandLine.StringVar(&passenv, "p", defaultPassenv, "environment var name of password")
	commandLine.BoolVar(&overwrite, "Y", overwrite, "overwrite output file")
	commandLine.Usage = func() { fmt.Print(encryptUsage) }
	commandLine.Parse(args)

	src, err := readFile(file)
	if err != nil {
		panic(err)
	}

	password, ok := os.LookupEnv(passenv)
	if !ok {
		panic(fmt.Sprintf("undefined environment variable '%s'", passenv))
	}

	data, err := encryptData(src, password)
	if err != nil {
		panic(err)
	}

	if err := writeFile(out, overwrite, magicHeader, data); err != nil {
		panic(err)
	}

	return 0
}
