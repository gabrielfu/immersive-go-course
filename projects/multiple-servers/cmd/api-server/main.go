package main

import (
	"flag"
	"fmt"
	"os"
	"servers/api"
)

func main() {
	dbUrl := os.Getenv("DATABASE_URL")
	port := flag.Int("port", 8081, "")
	flag.Parse()
	err := api.Run(dbUrl, *port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
