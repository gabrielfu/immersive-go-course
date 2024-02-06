package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net/http"
)

func main() {
	schedule := flag.String("schedule", "", "The schedule for the job")
	command := flag.String("command", "", "The command to run")
	port := flag.String("port", "8080", "The port to run the API server on")
	flag.Parse()

	r, err := http.Post(
		"http://localhost:"+*port+"/jobs",
		"application/json",
		bytes.NewBuffer([]byte(`{"schedule":"`+*schedule+`","command":"`+*command+`"}`)),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	log.Println("Status:", r.Status)
	log.Println("Header:", r.Header)
	body, _ := io.ReadAll(r.Body)
	log.Println("Body:", string(body))
}
