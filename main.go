package main

import (
	"fmt"
	"io"
	"os"

	"github.com/takumakei/envy/expandenv"
)

var VERSION = "0.0.1"

func main() {
	os.Exit(run())

	src := `${YOU:=${YOU1:=$LANG}} [$YOU] [$YOU1] YOU`
	fmt.Println(src)
	fmt.Println(expandenv.ExpandEnv(src))
}

func run() int {
	args := os.Args[1:]
	if len(args) == 0 {
		usage(os.Stdout)
		return 0
	}

	command := args[0]
	switch command {
	case "help", "hel", "he", "h", "-h", "--help":
		usage(os.Stdout)

	case "envy", "env", "en", "e", "exec", "exe", "ex":
		return execute(args[1:])

	case "encrypt", "encryp", "encry", "encr", "enc", "E":
		return encrypt(args[1:])

	case "decrypt", "decryp", "decry", "decr", "dec", "de", "d", "D":
		return decrypt(args[1:])

	case "version", "ver", "v":
		return version(args[1:])

	default:
		fmt.Fprintf(os.Stderr, "error: unknown command '%s'\n", command)
		usage(os.Stdout)
		return 1
	}

	return 0
}

func usage(out io.Writer) {
	fmt.Fprintf(
		out,
		`usage: envy [command] [options] [arguments]

  execute a program with enviroment variables

commands:
  help     (alias: hel, he, h, -h, --help)
  envy     (alias: env, en, e, exec, exe, ex)
  encrypt  (alias: enc, E)
  decrypt  (alias: dec, D)
  version  (alias: ver, v)
`,
	)
}

func version(args []string) int {
	fmt.Println("envy version " + VERSION)
	return 0
}
