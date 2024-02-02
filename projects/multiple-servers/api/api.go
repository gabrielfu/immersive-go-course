package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/rs/cors"
)

func createTable(ctx context.Context, conn *pgx.Conn) error {
	_, err := conn.Exec(ctx, `CREATE TABLE public.images (
		id serial NOT NULL,
		title text NOT NULL,
		url text NOT NULL,
		alt_text text,
		PRIMARY KEY (id)
	);`)
	return err
}

func initTable(ctx context.Context, conn *pgx.Conn) error {
	for _, image := range initialData {
		sql := fmt.Sprintf(
			`INSERT INTO public.images(title, url, alt_text) VALUES ('%s', '%s', '%s');`,
			image.Title,
			image.URL,
			image.AltText,
		)
		_, err := conn.Exec(ctx, sql)
		if err != nil {
			return err
		}
	}
	return nil
}

func dropTableIfExists(ctx context.Context, conn *pgx.Conn) error {
	_, err := conn.Exec(ctx, `DROP TABLE IF EXISTS public.images;`)
	return err
}

func marshalIndent(payload any, indent string) ([]byte, error) {
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

func handleGetImages(ctx context.Context, conn *pgx.Conn, indent string) ([]byte, error) {
	images, err := fetchImages(ctx, conn)
	if err != nil {
		return nil, err
	}
	b, err := marshalIndent(images, indent)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func handlePostImages(ctx context.Context, conn *pgx.Conn, body io.ReadCloser, indent string) ([]byte, error) {
	var image Image
	rb, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rb, &image)
	if err != nil {
		return nil, err
	}
	err = createImage(ctx, conn, image)
	if err != nil {
		return nil, err
	}
	b, err := marshalIndent(image, indent)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Run(dbUrl string, port int) error {
	log.Printf("dbUrl: %s\n", dbUrl)
	log.Printf("port: %d\n", port)

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		return fmt.Errorf("unable to connect to database %s: %v", dbUrl, err)
	}
	defer conn.Close(ctx)

	err = dropTableIfExists(ctx, conn)
	if err != nil {
		return fmt.Errorf("unable to drop table: %v", err)
	}
	err = createTable(ctx, conn)
	if err != nil {
		return fmt.Errorf("unable to create table: %v", err)
	}
	err = initTable(ctx, conn)
	if err != nil {
		return fmt.Errorf("unable to init table: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			indent := r.URL.Query().Get("indent")
			b, err := handleGetImages(ctx, conn, indent)
			if err != nil {
				fmt.Printf("failed to handle get images: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Add("Content-Type", "text/json")
			w.Write(b)
		case http.MethodPost:
			indent := r.URL.Query().Get("indent")
			b, err := handlePostImages(ctx, conn, r.Body, indent)
			if err != nil {
				fmt.Printf("failed to handle post images: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Add("Content-Type", "text/json")
			w.Write(b)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
	handler := cors.Default().Handler(mux)
	http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
	return nil
}
