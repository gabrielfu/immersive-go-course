package api

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

var initialData = []Image{
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

type Image struct {
	Title   string `json:"title"`
	AltText string `json:"alt_text"`
	URL     string `json:"url"`
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

func createImage(ctx context.Context, conn *pgx.Conn, image Image) error {
	sql := fmt.Sprintf(
		`INSERT INTO public.images(title, url, alt_text) VALUES ('%s', '%s', '%s');`,
		image.Title,
		image.URL,
		image.AltText,
	)
	_, err := conn.Exec(ctx, sql)
	return err
}
