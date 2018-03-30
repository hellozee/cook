package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	ps "Cook/parser"
)

func init() {
	//Initializing all the hash maps
	newfileTimings = make(map[string]string)
	oldfileTimings = make(map[string]string)
	fileList = make(map[string]string)
}

func main() {
	//Reading the Recipe File
	temp, err := ioutil.ReadFile("Recipe")

	if err != nil {
		//Missing Recipe File
		os.Stderr.WriteString("No sane Recipe File found.\nMake sure you have a Recipe file with proper syntax\n")
		return
	}

	Recipe := string(temp)

	//Parsing the Recipe File
	parser := ps.NewParser(Recipe)
	err = parser.Parse()
	checkErr(err)
	generateFileList(parser, parser.CompilerDetails.Start)
	var jsonData []byte

	if _, err := os.Stat("Cooking/details.json"); err == nil {

		//Reading the details.json which contains the file names against their generated timestamps
		jsonFile, err := os.Open("Cooking/details.json")
		defer jsonFile.Close()

		bytes, _ := ioutil.ReadAll(jsonFile)
		err = json.Unmarshal(bytes, &hashJSONold)

		if err != nil {
			//Someone has tampered with the JSON file
			os.Stderr.WriteString("Error parsing Cooking/details.json\nPlease run the program again\n")
			os.Remove("Cooking/details.json")
			return
		}

		structToMap(hashJSONold.Body.Entity)

		compareAndCompile(parser)

	} else {
		_ = os.Mkdir("Cooking", 0777)
		generateList()
		compileFirst(parser)
	}

	fmt.Println("All files Compiled...")

	jsonData, err = json.MarshalIndent(hashJSONnew, "", " ")
	checkErr(err)

	file, err := os.OpenFile("Cooking/details.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0777)
	checkErr(err)
	defer file.Close()

	_, err = file.Write(jsonData)
	checkErr(err)

	linkAll(parser)

}
