package main

import (
	"os"

	"github.com/gocarina/gocsv"
)

type Input struct {
	Url string `csv:"url"`
}

type Output struct {
	Url    string `csv:"url"`
	Input  string `csv:"input"`
	Output string `csv:"output"`
	S3Url  string `csv:"s3url"`
}

func ReadInputs(filepath string) ([]Input, error) {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var inputs []Input
	if err := gocsv.UnmarshalFile(f, &inputs); err != nil {
		panic(err)
	}
	return inputs, nil
}

func WriteOutputs(filepath string, outputs []Output) error {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	err = gocsv.MarshalFile(&outputs, f)
	if err != nil {
		panic(err)
	}
	return nil
}
