package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func structToMap(parsedStruct []entity) {

	for _, item := range parsedStruct {
		oldfileTimings[item.File] = item.Hash
	}

}

func generateList() {

	for _, value := range fileList {
		file, err := os.Stat(value)
		checkErr(err)
		t := file.ModTime()
		hash, err := hashTime(t.String())
		checkErr(err)

		newfileTimings[value] = hash
		hashJSONnew.Body.Entity = append(hashJSONnew.Body.Entity, entity{File: value, Hash: hash})
	}
}

func compileFirst() {
	//Iteratively generate .o files

	for tag, file := range fileList {
		fmt.Println("Compiling " + file)
		cmd := exec.Command(compilerDetails.binary, "-c", file, "-o", "Cooking/"+tag+".o")
		checkCommand(cmd)
		tagList = append(tagList, "Cooking/"+tag+".o")
	}
}

func compareAndCompile() {
	//Compare the file hash with current hash if do not match generate .o file
	//also replace the current hash with the new hash

	for key, value := range fileList {

		file, err := os.Stat(value)

		checkErr(err)
		t := file.ModTime()

		if !checkTimeStamp(t.String(), oldfileTimings[value]) {
			fmt.Println("Compiling " + value)
			cmd := exec.Command(compilerDetails.binary, "-c", value, "-o", "Cooking/"+key+".o")
			checkCommand(cmd)

			oldfileTimings[value], err = hashTime(t.String())
			checkErr(err)
		}

		hashJSONnew.Body.Entity = append(hashJSONnew.Body.Entity, entity{File: value, Hash: oldfileTimings[value]})
		tagList = append(tagList, "Cooking/"+key+".o")
	}
}

func linkAll() {

	//Compile all the generated .o files under the Cooking directory
	fmt.Println("Linking files..")
	args := []string{compilerDetails.binary, "-o", compilerDetails.name, compilerDetails.includes, compilerDetails.otherFlags}
	for _, tag := range tagList {
		args = append(args, tag)
	}
	args = append(args, compilerDetails.ldFlags)
	cmd := exec.Command(os.Getenv("SHELL"), "-c", strings.Join(args, " "))
	checkCommand(cmd)
}
