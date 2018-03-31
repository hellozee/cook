package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	ps "github.com/hellozee/cook/parser"
)

func structToMap(parsedStruct []entity) {
	for _, item := range parsedStruct {
		oldfileTimings[item.File] = item.Hash
	}
}

func generateList() {
	for _, value := range fileList {
		file, err := ioutil.ReadFile(value)
		checkErr(err)
		hash := hashFile(file)
		newfileTimings[value] = hash
		hashJSONnew.Body.Entity = append(hashJSONnew.Body.Entity,
			entity{File: value, Hash: hash})
	}
}

func compileFirst(par ps.Parser) {
	//Iteratively generate .o files

	for key, value := range fileList {
		if *quietFlag == false {
			fmt.Println("Compiling " + value)
		}
		cmd := exec.Command(par.CompilerDetails.Binary, "-c", value,
			"-o", "Cooking/"+key+".o")
		checkCommand(cmd)
	}
}

func compareAndCompile(par ps.Parser) {
	//Compare the file hash with current hash if do not match generate .o file
	//also replace the current hash with the new hash

	for key, value := range fileList {

		file, err := ioutil.ReadFile(value)
		checkErr(err)

		if !checkHash(file, oldfileTimings[value]) {
			if *quietFlag == false {
				fmt.Println("Compiling " + value)
			}
			cmd := exec.Command(par.CompilerDetails.Binary, "-c", value,
				"-o", "Cooking/"+key+".o")
			checkCommand(cmd)

			oldfileTimings[value] = hashFile(file)
		}

		hashJSONnew.Body.Entity = append(hashJSONnew.Body.Entity,
			entity{File: value, Hash: oldfileTimings[value]})
	}
}

func linkAll(par ps.Parser) {

	//Compile all the generated .o files under the Cooking directory
	if *quietFlag == false {
		fmt.Println("Linking files..")
	}
	args := []string{par.CompilerDetails.Binary, "-o", par.CompilerDetails.Name,
		par.CompilerDetails.Includes, par.CompilerDetails.OtherFlags,
		"Cooking/*.o", par.CompilerDetails.LdFlags}
	cmd := exec.Command(os.Getenv("SHELL"), "-c", strings.Join(args, " "))
	checkCommand(cmd)
}
