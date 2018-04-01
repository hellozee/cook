package main

import (
	"flag"
	"fmt"
	"os"

	mg "github.com/hellozee/cook/manager"
	ps "github.com/hellozee/cook/parser"
	wk "github.com/hellozee/cook/worker"
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

	if err != nil {
		fmt.Println("Error Parsing the Recipe File :")
		fmt.Println(err.Error())
		return
	}

	worker := wk.Worker{
		QuietFlag:   *quietFlag,
		VerboseFlag: *verboseFlag,
	}

	if *cleanFlag == true {
		os.RemoveAll("Cooking/")
		os.Remove(parser.CompilerDetails.Name)
		return
	}

	err = manager.GenerateFileList(parser, parser.CompilerDetails.Start)

	if _, err := os.Stat("Cooking/details.json"); err == nil {

		err = manager.ReadDetails()

		if err != nil {
			fmt.Println("Unable to read details.json :")
			fmt.Println(err.Error())
			return
		}

		worker.CompareAndCompile(parser, &manager)

	} else {
		_ = os.Mkdir("Cooking", 0755)
		err = manager.GenerateList()
		if err != nil {
			fmt.Println("Error Generating the file list :")
			fmt.Println(err.Error())
			return
		}
		worker.CompileFirst(parser, manager)
	}

	if *quietFlag == false {
		fmt.Println("All files Compiled...")
	}

	err = manager.WriteDetails()

	if err != nil {
		fmt.Println("Unable to write details.json :")
		fmt.Println(err.Error())
		return
	}

	err = worker.Link(parser)

	if err != nil {
		fmt.Println("Unable to Link Files:")
		fmt.Println(err.Error())
		return
	}

}
