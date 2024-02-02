package main

import (
	"flag"
	"servers/static"
)

func main() {
	path := flag.String("path", "assets", "")
	port := flag.Int("port", 8082, "")
	flag.Parse()
	static.Run(*path, *port)
}
