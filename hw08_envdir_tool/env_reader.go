package main

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

var ErrEmptyDir = errors.New("dir must not be empty")

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	if dir == "" {
		return nil, ErrEmptyDir
	}
	dirEntry, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	env := make(Environment)
	for _, file := range dirEntry {
		if file.IsDir() {
			continue
		}
		if strings.Contains(file.Name(), "=") {
			continue
		}

		fPath := filepath.Join(dir, file.Name())

		stat, err := os.Stat(fPath)
		if err != nil {
			return nil, err
		}

		fName := file.Name()

		if stat.Size() == 0 {
			env[fName] = EnvValue{
				Value:      "",
				NeedRemove: true,
			}
			continue
		}

		file, err := os.Open(fPath)
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(file)
		if scanner.Scan() {
			str := scanner.Text()
			str = strings.TrimRight(str, " \t")
			str = string(bytes.ReplaceAll([]byte(str), []byte("\x00"), []byte("\n")))
			env[fName] = EnvValue{
				Value:      str,
				NeedRemove: false,
			}
		}
	}

	return env, nil
}
