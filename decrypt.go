package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
)

var decryptUsage = `usage: envy decrypt [options]

  decrypt the file

options:
  -h           show this message
  -f <file>    input filename to decrypt
               (default "-" stdin)
  -p <env>     enviromnent variable name of the password
               (default "ENVY_PASSWORD")
  -o <file>    output filename
               (default "-" stdout)
  -Y           allow to overwrite the output file
`

func decrypt(args []string) int {
	var out string
	var file string
	var passenv string
	var overwrite bool
	commandLine := flag.NewFlagSet("decrypt", flag.ExitOnError)
	commandLine.StringVar(&out, "o", out, "output file")
	commandLine.StringVar(&file, "f", file, "input file")
	commandLine.StringVar(&passenv, "p", defaultPassenv, "enviroment var name of password")
	commandLine.BoolVar(&overwrite, "Y", overwrite, "overwrite output file")
	commandLine.Usage = func() { fmt.Print(decryptUsage) }
	commandLine.Parse(args)

	src, err := readFile(file)
	if err != nil {
		panic(err)
	}

	if !bytes.HasPrefix(src, magicHeader) {
		panic("bad magic")
	}
	src = src[len(magicHeader):]

	password, ok := os.LookupEnv(passenv)
	if !ok {
		panic(fmt.Sprintf("undefined environment variable '%s'", passenv))
	}

	data, err := decryptData(src, password)
	if err != nil {
		panic(err)
	}

	if err := writeFile(out, overwrite, data); err != nil {
		panic(err)
	}

	return 0
}
