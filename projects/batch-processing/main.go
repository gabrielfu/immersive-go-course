package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/gographics/imagick.v2/imagick"
)

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

	// Image processing
	var outputs []Output
	for _, input := range inputs {
		log.Println("downloading image from", input.Url)
		src, err := DownloadImage(input.Url)
		if err != nil {
			log.Printf("error: %v\n", err)
			continue
		}

		log.Println("processing image", src)

		// Generate temp file name
		dest, err := NewTempFileName()
		if err != nil {
			log.Printf("error: %v\n", err)
			continue
		}

		// Do the conversion
		err = c.Grayscale(src, dest)
		if err != nil {
			log.Printf("error: %v\n", err)
			continue
		}
		log.Printf("saved to %s\n", dest)

		// Upload to S3
		UploadFile(s3, bucket, filepath.Base(dest), dest)

		// Add to the outputs
		outputs = append(outputs, Output{
			Url:    input.Url,
			Input:  src,
			Output: dest,
			S3Url:  "https://" + bucket + ".s3.amazonaws.com/" + filepath.Base(dest),
		})
	}
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
