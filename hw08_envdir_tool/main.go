package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go-envdir /path/to/directory cmd arg=value...")
		os.Exit(1)
	}

	dir := os.Args[1]
	cmd := os.Args[2:]

	env, err := ReadDir(dir)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	rc := RunCmd(cmd, env)

	os.Exit(rc)
}
