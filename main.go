package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	ps "github.com/hellozee/cook/parser"
)

func init() {
	//Initializing all the hash maps
	newfileTimings = make(map[string]string)
	oldfileTimings = make(map[string]string)
	fileList = make(map[string]string)
}

var quietFlag = flag.Bool("quiet", false, "To not show any output")
var cleanFlag = flag.Bool("clean", false, "To clean the cached data")
var helpFlag = flag.Bool("help", false, "To show this help message")
var verboseFlag = flag.Bool("verbose", false, "To increase the level of verbosity")

func main() {

	flag.Parse()

	help := `
	Usage: cook [OPTIONS]

	--help:
		To show this help message

	--quiet:
		To not show any output

	--clean:
		To clean the cached data and perform a clean build

	--verbose:
		To increase the verbosity level
	`
	if *helpFlag == true {
		fmt.Println(help)
		return
	}

	//Reading the Recipe File
	temp, err := ioutil.ReadFile("Recipe")

	if err != nil {
		//Missing Recipe File
		os.Stderr.WriteString("No sane Recipe File found.\n" +
			"Make sure you have a Recipe file with proper syntax\n")
		return
	}

	Recipe := string(temp)

	//Parsing the Recipe File
	parser := ps.NewParser(Recipe)
	err = parser.Parse()
	checkErr(err)

	if *cleanFlag == true {
		os.RemoveAll("Cooking/")
		os.Remove(parser.CompilerDetails.Name)
		return
	}

	generateFileList(parser, parser.CompilerDetails.Start)
	var jsonData []byte

	if _, err := os.Stat("Cooking/details.json"); err == nil {

		//Reading the details.json which contains the file names
		//against their generated timestamps
		jsonFile, err := os.Open("Cooking/details.json")
		defer jsonFile.Close()

		bytes, _ := ioutil.ReadAll(jsonFile)
		err = json.Unmarshal(bytes, &hashJSONold)

		if err != nil {
			//Someone has tampered with the JSON file
			os.Stderr.WriteString("Error parsing Cooking/details.json\n" +
				"Please run the program again\n")
			os.Remove("Cooking/details.json")
			return
		}

		structToMap(hashJSONold.Body.Entity)

		compareAndCompile(parser)

	} else {
		_ = os.Mkdir("Cooking", 0755)
		generateList()
		compileFirst(parser)
	}

	if *quietFlag == false {
		fmt.Println("All files Compiled...")
	}

	jsonData, err = json.MarshalIndent(hashJSONnew, "", " ")
	checkErr(err)

	file, err := os.OpenFile("Cooking/details.json",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	checkErr(err)
	defer file.Close()

	_, err = file.Write(jsonData)
	checkErr(err)

	linkAll(parser)

}
