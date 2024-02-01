package main

import (
	"encoding/csv"
	"fmt"
	"io"
)

//type Parameter interface {
//	any | []byte
//}

func pipelineStage[IN any, OUT any](input <-chan IN, output chan<- OUT, process func(IN) OUT) {
	defer close(output)
	for data := range input {
		output <- process(data)
	}
}

func asynchronousPipeline(input *csv.Reader) {
	fmt.Println("--Asynchronous pipeline----")
	parseInputCh := make(chan []string)
	convertInputCh := make(chan Record)
	encodeInputCh := make(chan Record)

	// read the output of the pipeline from this channel
	outputCh := make(chan []byte)
	// channel to wait for the printing of the final result
	done := make(chan struct{})

	//	start pipeline stages and connect them

	go pipelineStage(parseInputCh, convertInputCh, parse)
	go pipelineStage(convertInputCh, encodeInputCh, convert)
	//go pipelineStage(encodeInputCh, done, encode)

	// start go routine to read pipeline output
	go func() {
		for data := range outputCh {
			fmt.Println(string(data))
		}
		close(done)
	}()

	// ignore the 1st row
	input.Read()

	for {
		rec, err := input.Read()
		if err == io.EOF {
			close(parseInputCh)
			break
		}
		if err != nil {
			panic(err)
		}
		// send input to pipeline
		parseInputCh <- rec
	}

	// wait until the last output is finished
	<-done
}
