package main

import (
	"flag"
	"fmt"
	"os"
	"parser"
)

// expecting  -i[file] -o[file2] -c[file3]
func main() {

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

	data, err := parser.ReadJSON(input)
	if err != nil {
		panic(err)
	}
	for _, jsonObj := range data.([]map[string]interface{}) {
		parser.Filter(jsonObj, jsonConfig.([]map[string]interface{})[0])
	}
	output, err := os.Create(outFName)
	if err != nil {
		fmt.Printf("failed to create new file - %v ", err)
	}
	defer output.Close()

	if err := parser.WriteJSON(data, output); err != nil {
		fmt.Printf("couldn't write json: %v\n", err)
	}
	fmt.Println(data)

}
