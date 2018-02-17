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

func generateList(tag string) {
	//Generate Cooking/details.json
	parameters := fileDetails[tag]

	file, err := os.Stat(parameters.file)
	checkErr(err)

	t := file.ModTime()
	hash, err := hashTime(t.String())
	checkErr(err)

	newfileTimings[parameters.file] = hash
	hashJSONnew.Body.Entity = append(hashJSONnew.Body.Entity, entity{File: parameters.file, Hash: hash})

	if parameters.deps == nil {
		return
	}

	for _, name := range parameters.deps {
		generateList(name)
	}
}

func compileFirst(tag string) {
	//Recursively generate .o files

	parameters := fileDetails[tag]
	fmt.Println("Compiling " + parameters.file)
	cmd := exec.Command(compilerDetails.binary, "-c", parameters.file, "-o", "Cooking/"+tag+".o")
	checkCommand(cmd)
	tagList = append(tagList, "Cooking/"+tag+".o")

	for _, name := range parameters.deps {
		compileFirst(name)
	}
}

func compareAndCompile(tag string) {
	//Compare the file hash with current hash if do not match generate .o file
	//also replace the current hash with the new hash

	parameters := fileDetails[tag]

	file, err := os.Stat(parameters.file)

	checkErr(err)
	t := file.ModTime()

	if !checkTimeStamp(t.String(), oldfileTimings[parameters.file]) {
		fmt.Println("Compiling " + parameters.file)
		cmd := exec.Command(compilerDetails.binary, "-c", parameters.file, "-o", "Cooking/"+tag+".o")
		checkCommand(cmd)

		oldfileTimings[parameters.file], err = hashTime(t.String())
		checkErr(err)
	}

	hashJSONnew.Body.Entity = append(hashJSONnew.Body.Entity, entity{File: parameters.file, Hash: oldfileTimings[parameters.file]})
	tagList = append(tagList, "Cooking/"+tag+".o")
	for _, name := range parameters.deps {
		compareAndCompile(name)
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
