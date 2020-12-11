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
		fmt.Printf("faild to open file- %s:\n", configFName)
		panic(err)
	}
	defer configFilePtr.Close()

	jsonConfig := parser.ReadJSON(configFilePtr)
	input, err := os.Open(inputFName)
	if err != nil {
		fmt.Printf("faild to open file- %s:\n", inputFName)
		panic(err)
	}
	defer input.Close()

	data := parser.ReadJSON(input)
	for _, jsonObj := range data.([]map[string]interface{}) {
		parser.Filter(jsonObj, jsonConfig.([]map[string]interface{})[0])

	}
	output, err := os.Create(outFName)
	if err != nil {
		fmt.Printf("failed to create new file - %s ", outFName)
	}
	defer output.Close()

	parser.WriteJSON(data, output)
	fmt.Println(data)

}
