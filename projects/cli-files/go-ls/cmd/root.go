package cmd

import (
	"fmt"
	"os"
)

func ls(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		fmt.Println(path)
		return nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fmt.Println(entry.Name())
	}
	return nil
}

func Execute() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	if err := ls(path); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
