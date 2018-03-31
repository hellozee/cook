package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	mg "github.com/hellozee/cook/manager"
	ps "github.com/hellozee/cook/parser"
)

func compileFirst(par ps.Parser, man mg.Manager) {
	//Iteratively generate .o files

	for key, value := range man.FileList {
		if *quietFlag == false {
			fmt.Println("Compiling " + value)
		}
		cmd := exec.Command(par.CompilerDetails.Binary, "-c", value,
			"-o", "Cooking/"+key+".o")
		checkCommand(cmd)
	}
}

func compareAndCompile(par ps.Parser, man mg.Manager) {
	//Compare the file hash with current hash if do not match generate .o file
	//also replace the current hash with the new hash

	for key, value := range man.FileList {

		file, err := ioutil.ReadFile(key)

		checkErr(err)

		if !mg.CheckHash(file, man.OldFileTimings[value]) {

			if *quietFlag == false {
				fmt.Println("Compiling " + value)
			}
			cmd := exec.Command(par.CompilerDetails.Binary, "-c", value,
				"-o", "Cooking/"+key+".o")
			checkCommand(cmd)
			man.OldFileTimings[value] = mg.HashFile(file)

		}

		man.HashJSONnew.Body.Entity = append(man.HashJSONnew.Body.Entity,
			mg.Entity{File: value, Hash: man.OldFileTimings[value]})
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
