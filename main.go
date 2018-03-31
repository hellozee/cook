package main

import (
	"flag"
	"fmt"
	"os"

	mg "github.com/hellozee/cook/manager"
	ps "github.com/hellozee/cook/parser"
)

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
		To clean the cached data

	--verbose:
		To increase the verbosity level
		
	`
	if *helpFlag == true {
		fmt.Println(help)
		return
	}

	//Reading the Recipe File
	manager, err := mg.NewManager()

	Recipe := string(manager.FileData)

	//Parsing the Recipe File
	parser := ps.NewParser(Recipe)
	err = parser.Parse()
	checkErr(err)

	if *cleanFlag == true {
		os.RemoveAll("Cooking/")
		os.Remove(parser.CompilerDetails.Name)
		return
	}

	err = manager.GenerateFileList(parser, parser.CompilerDetails.Start)

	if _, err := os.Stat("Cooking/details.json"); err == nil {

		err = manager.ReadDetails()

		compareAndCompile(parser, manager)

	} else {
		_ = os.Mkdir("Cooking", 0755)
		err = manager.GenerateList()
		compileFirst(parser, manager)
	}

	if *quietFlag == false {
		fmt.Println("All files Compiled...")
	}

	err = manager.WriteDetails()

	linkAll(parser)

}
