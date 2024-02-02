package static

import (
	"fmt"
	"log"
	"net/http"
)

func Run(path string, port int) {
	log.Printf("path: %s\n", path)
	log.Printf("port: %d\n", port)

	http.Handle("/", http.FileServer(http.Dir(path)))
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
