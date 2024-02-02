package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gopkg.in/gographics/imagick.v2/imagick"
)

func Pipeline(s3 *s3.Client, bucket string, url string, c *Converter) (*Output, error) {
	log.Printf("[%s] downloading image\n", url)
	src, err := DownloadImage(url)
	if err != nil {
		return nil, err
	}
	log.Printf("[%s] saved to %s\n", url, src)

	// Generate temp file name
	dest, err := NewTempFileName()
	if err != nil {
		return nil, err
	}

	// Do the conversion
	log.Printf("[%s] processing image to %s\n", url, dest)
	err = c.Grayscale(src, dest)
	if err != nil {
		return nil, err
	}

	// Upload to S3
	log.Printf("[%s] uploading image to s3\n", url)
	UploadFile(s3, bucket, filepath.Base(dest), dest)

	// Add to the outputs
	return &Output{
		Url:    url,
		Input:  src,
		Output: dest,
		S3Url:  "https://" + bucket + ".s3.amazonaws.com/" + filepath.Base(dest),
	}, nil
}

func main() {
	// Accept --input and --output arguments for the images
	inputFilepath := flag.String("input", "", "A path to the csv file containing image urls")
	outputFilepath := flag.String("output", "", "A path to where the output csv should be written")
	flag.Parse()

	s3 := NewS3()
	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		log.Println("S3_BUCKET environment variable not set")
		os.Exit(1)
	}

	// Log what we're going to do
	log.Printf("processing: %q to %q\n", *inputFilepath, *outputFilepath)
	inputs, err := ReadInputs(*inputFilepath)
	if err != nil {
		log.Printf("error: %v\n", err)
		os.Exit(1)
	}

	// Set up imagemagick
	imagick.Initialize()
	defer imagick.Terminate()

	// Build a Converter struct that will use imagick
	c := &Converter{
		cmd: imagick.ConvertImageCommand,
	}

	wg := sync.WaitGroup{}

	// Image processing
	var outputs []Output
	for _, input := range inputs {
		wg.Add(1)
		go func(url string) {
			output, err := Pipeline(s3, bucket, url, c)
			if err != nil {
				log.Printf("[%s] error: %v\n", url, err)
			} else {
				outputs = append(outputs, *output)
			}
			wg.Done()
		}(input.Url)
	}
	wg.Wait()
	log.Println("done processing")
	log.Println(outputs)

	// Write the outputs to a file
	log.Printf("writing outputs to %q\n", *outputFilepath)
	err = WriteOutputs(*outputFilepath, outputs)
	if err != nil {
		log.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
