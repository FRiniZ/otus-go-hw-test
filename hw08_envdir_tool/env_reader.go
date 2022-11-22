package main

import (
	"bufio"
	"bytes"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func getEnvValue(filename string) (*EnvValue, error) {
	eV := &EnvValue{
		Value:      "",
		NeedRemove: true,
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	if line := scanner.Text(); line != "" {
		line := scanner.Text()
		line = strings.TrimRight(line, " \t")
		line = string(bytes.ReplaceAll([]byte(line), []byte("\x00"), []byte("\n")))
		if line != "" {
			eV.Value = line
			eV.NeedRemove = false
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return eV, nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	Env := Environment{}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, s := range files {
		if s.IsDir() {
			continue
		}

		t, err := getEnvValue(dir + "/" + s.Name())
		if err != nil {
			return nil, err
		}
		Env[s.Name()] = *t
	}

	return Env, nil
}
