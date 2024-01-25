package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
)

type Image struct {
	Title   string `json:"title"`
	AltText string `json:"alt_text"`
	URL     string `json:"url"`
}

var data = []Image{
	{
		Title:   "Sunset",
		AltText: "Clouds at sunset",
		URL:     "https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80",
	},
	{
		Title:   "Mountain",
		AltText: "A mountain at sunset",
		URL:     "https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80",
	},
}

func MarshalIndent(payload any, indent string) ([]byte, error) {
	if indent == "" {
		return json.Marshal(payload)
	}

	num, err := strconv.Atoi(indent)
	if err != nil {
		return nil, err
	}
	indentStr := strings.Repeat(" ", num)
	return json.MarshalIndent(payload, "", indentStr)
}

func fetchImages(ctx context.Context, conn *pgx.Conn) ([]Image, error) {
	var images []Image
	rows, err := conn.Query(ctx, "SELECT title, url, alt_text FROM public.images;")
	if err != nil {
		return nil, err
	}

	var title, url, altText string
	for rows.Next() {
		err = rows.Scan(&title, &url, &altText)
		if err != nil {
			return nil, err
		}
		images = append(images, Image{Title: title, URL: url, AltText: altText})
	}
	return images, nil
}

func main() {
	dbUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		fmt.Fprintln(os.Stderr, "'DATABASE_URL' is not set")
		os.Exit(1)
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, `CREATE TABLE public.images (
		id serial NOT NULL,
		title text NOT NULL,
		url text NOT NULL,
		alt_text text,
		PRIMARY KEY (id)
	);`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create table: %v\n", err)
		os.Exit(1)
	}

	for _, image := range data {
		sql := fmt.Sprintf(
			`INSERT INTO public.images(title, url, alt_text) VALUES ('%s', '%s', '%s');`,
			image.Title,
			image.URL,
			image.AltText,
		)
		_, err = conn.Exec(ctx, sql)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to insert row: %v\n", err)
			os.Exit(1)
		}
	}

	http.HandleFunc("/images.json", func(w http.ResponseWriter, r *http.Request) {
		images, err := fetchImages(ctx, conn)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		indent := r.URL.Query().Get("indent")
		b, err := MarshalIndent(images, indent)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "text/json")
		w.Write(b)
	})

	http.ListenAndServe(":8080", nil)
}

// To run:
// > docker run -e POSTGRES_DB=db -e POSTGRES_USER=username -e POSTGRES_PASSWORD=password -p 5432:5432 postgres:alpine3.18
// > DATABASE_URL='postgres://username:password@localhost:5432/db' go run .
//
// To request:
// > curl 'http://localhost:8080/images.json?indent=2' -i
