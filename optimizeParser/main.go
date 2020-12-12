package main

import (
	"flag"
	"fmt"
	"optimizeParser"
	"os"
	"parser"
	"sync"
)

// expecting  -i[file] -o[file2] -c[file3]
func main() {
	var wg sync.WaitGroup
	wg.Add(3)
	var inputFName, outFName, configFName string
	if len(os.Args) < 4 {
		panic("command args missing")
	}
	flag.StringVar(&inputFName, "i", "", "json file to parse ")
	flag.StringVar(&configFName, "c", "", "configuration file")
	flag.StringVar(&outFName, "o", "", "output file name")
	flag.Parse()

	configFilePtr, err := os.Open(configFName)
	if err != nil {
		panic(err)
	}
	defer configFilePtr.Close()
	jsonConfig, err := parser.ReadJSON(configFilePtr)
	input, err := os.Open(inputFName)
	if err != nil {
		panic(err)
	}
	defer input.Close()

	writerCh := make(chan map[string]interface{})
	ReaderTofilterCh := make(chan map[string]interface{}, 2)

	go func() {
		optimizeParser.ReadJSONAndSendOverChannel(input, ReaderTofilterCh)
		wg.Done()
	}()
	go func() {
		optimizeParser.ProcessjsonUsingConfig(jsonConfig[0], ReaderTofilterCh, writerCh)
		wg.Done()
	}()
	output, err := os.Create(outFName)
	if err != nil {
		fmt.Printf("failed to create new file - %v ", err)
	}
	defer output.Close()
	go func() {
		optimizeParser.GetDataFromChannelAndWriteJSON(writerCh, output)
		wg.Done()
	}()

	wg.Wait()
}
