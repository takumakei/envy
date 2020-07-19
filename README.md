ENVY = ENV + YAML
======================================================================

`/usr/bin/env` with additional features.

- load enviroment variables defined in YAML files encrypted optionally
- encrypt and decrypt files
  - AES-256 with Galois Counter Mode
- unset variables
- `${var:-expr}` and `${var:=expr}`


```help
usage: envy [command] [options] [arguments]

  execute a program with enviroment variables

commands:
  help     (alias: hel, he, h, -h, --help)
  env      (alias: en, e, exec, exe, ex)
  encrypt  (alias: enc, E)
  decrypt  (alias: dec, D)
```


env command
----------------------------------------------------------------------

```help
usage: envy env [options] [--] [command [args...]]

  execute command with specified enviromnent variables set

options:
  Each option can be used multiple times.
  The order of the options is sensitive.
  Options that you use later overwrite options that you used earlier.

  -h, --help            show this help.
  -p, --password <env>  decrypt yaml files loaded after this option
                        using the password read from the environment
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

yaml file styles:
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
```


encrypt command
----------------------------------------------------------------------

```help
usage: envy decrypt [options]

  decrypt the file

options:
  -h           show this message
  -f <file>    input filename to decrypt
  -p <env>     enviromnent variable name of the password
               (default "ENVY_PASSWORD")
  -o <file>    output filename
  -Y           allow overwrite the output file
```


decrypt command
----------------------------------------------------------------------

```help
usage: envy decrypt [options]

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
```


LICENSE
----------------------------------------------------------------------

MIT License

Copyright (c) 2020 takumakei
