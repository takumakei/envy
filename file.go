package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/natefinch/atomic"
)

var magicHeader = []byte("ENVY")

func readFile(file string) ([]byte, error) {
	if file == "" || file == "-" {
		return ioutil.ReadAll(os.Stdin)
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func writeFile(file string, overwrite bool, data ...[]byte) error {
	if file == "" || file == "-" {
		_, err := io.Copy(os.Stdout, NewReader(data...))
		return err
	}

	if !overwrite {
		_, err := os.Stat(file)
		if err == nil {
			return fmt.Errorf("file exists '%s'", file)
		}
		if !os.IsNotExist(err) {
			return err
		}
	}

	return atomic.WriteFile(file, NewReader(data...))
}
