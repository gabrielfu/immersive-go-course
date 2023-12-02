package cmd

import (
	"fmt"
	"os"
)

func cat(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return fmt.Errorf("%s is a directory", path)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read %s: %w", path, err)
	}

	os.Stdout.Write(content)
	return nil
}

func Execute() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	if err := cat(path); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
