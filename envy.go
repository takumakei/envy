package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
	"syscall"

	"github.com/takumakei/envy/expandenv"
	"gopkg.in/yaml.v2"
)

var executeUsage = `usage: envy envy [options] [--] [command [args...]]

  execute command with specified enviromnent variables set

options:
  Each option can be used multiple times.
  The order of the options is sensitive.
  Options that you use later overwrite options that you used earlier.

  -h, --help            show this help.
  -p, --password <env>  decrypt yaml files loaded after this option.
                        using the password read form the environment
                        variable.
                        use a empty string for <env> in case of that
                        yaml file is not encrypted.
  -f, --file <file>     read environment variables from the yaml file.
  <key=value>           setenv key=value

expansion of a variable in a value:
  Values can contain ${var:-expr} or ${var:=expr}.
  The expr is used in case of var is undefined.
  The var environment variable is also set in case of using ':='.
  The expr can also contain ${var:-expr} or ${var:=expr} recursively.

yaml file format:
  Files can be written in one of three styles.
  Environment variables can be unset with a null value using a map style
  or a list of map styles.
  There is no way to unset a variable in the list of <key=value> style.
  Elements in a map are sorted by its key.
  Use the list of maps style to specify the order.

  * a map style

    KEY_0: VALUE_0
    KEY_1: VALUE_1

  * a list of maps style

    - KEY_0: VALUE_0
    - KEY_1: VALUE_1
    - KEY_2: # unset KEY_2 because of the value is null

  * a list of <key=value> style

    - KEY_0=VALUE_0
    - KEY_1=VALUE_1
`

func execute(args []string) int {
	var password string
	for n := len(args); n > 0; n = len(args) {
		arg := args[0]
		if arg == "--" {
			args = args[1:]
			break
		}

		if arg == "-h" || arg == "--help" {
			fmt.Print(executeUsage)
			return 0
		}

		if arg == "-f" || arg == "--file" {
			if n <= 1 {
				panic("no arguments")
			}
			if err := parseEnvYaml(args[1], password); err != nil {
				panic(err)
			}
			args = args[2:]
			continue
		}

		if arg == "-p" || arg == "--password" {
			if n <= 1 {
				panic("no password <env>")
			}
			password = args[1]
			args = args[2:]
			continue
		}

		if strings.HasPrefix(arg, "-f=") || strings.HasPrefix(arg, "--file=") {
			m := strings.SplitN(arg, "=", 2)
			if err := parseEnvYaml(m[1], password); err != nil {
				panic(err)
			}
			args = args[1:]
			continue
		}

		if strings.HasPrefix(arg, "-") {
			panic("unknown option")
		}

		if m := strings.SplitN(arg, "=", 2); len(m) == 2 {
			if err := os.Setenv(m[0], expandenv.ExpandEnv(m[1])); err != nil {
				panic(err)
			}
			args = args[1:]
			continue
		}

		break
	}

	if len(args) == 0 {
		printEnv()
	} else {
		execArgs(args)
	}
	return 0
}

func parseEnvYaml(file, password string) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	if password != "" {
		p, ok := os.LookupEnv(password)
		if !ok {
			panic(fmt.Sprintf("password env '%s' not defined", password))
		}
		if p != "" {
			if !bytes.HasPrefix(b, magicHeader) {
				panic("not envy file")
			}
			b = b[len(magicHeader):]
			b, err = decryptData(b, p)
			if err != nil {
				panic(err)
			}
		}
	}

	var y0 []string
	err = yaml.Unmarshal(b, &y0)
	if err == nil {
		return setenvList(y0)
	}

	var y1 map[string]*string
	err = yaml.Unmarshal(b, &y1)
	if err == nil {
		return setenvMap(y1)
	}

	var y2 []map[string]*string
	err = yaml.Unmarshal(b, &y2)
	if err == nil {
		return setenvMapList(y2)
	}

	return err
}

func setenvList(list []string) error {
	for _, e := range list {
		m := strings.SplitN(e, "=", 2)
		if len(m) != 2 {
			return fmt.Errorf("invalid format '%s'", e)
		}
		if err := os.Setenv(m[0], expandenv.ExpandEnv(m[1])); err != nil {
			return err
		}
	}
	return nil
}

func setenvMap(m map[string]*string) error {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := m[k]
		if v == nil {
			if err := os.Unsetenv(k); err != nil {
				return err
			}
		} else if err := os.Setenv(k, expandenv.ExpandEnv(*v)); err != nil {
			return err
		}
	}
	return nil
}

func setenvMapList(m []map[string]*string) error {
	for _, e := range m {
		if err := setenvMap(e); err != nil {
			return err
		}
	}
	return nil
}

func printEnv() {
	list := os.Environ()
	for _, e := range list {
		fmt.Println(e)
	}
}

func execArgs(args []string) {
	binary, err := exec.LookPath(args[0])
	if err != nil {
		panic(err)
	}

	env := os.Environ()

	err = syscall.Exec(binary, args, env)
	if err != nil {
		panic(err)
	}
}
