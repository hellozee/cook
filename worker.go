package main

import (
	"os"
)

func structToMap(parsedStruct []entity) {

	for _, item := range parsedStruct {
		oldfileTimings[item.File] = item.Hash
	}

}

func generateList(tag string) {

	parameters := fileDetails[tag]

	file, err := os.Stat(parameters.file)
	//fmt.Println(tag)
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

	//exec.Command(compilerDetails.binary, "-c", parameters.file, "-o", "Cooking/"+tag+".o")

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
		//exec.Command(compilerDetails.binary, "-c", parameters.file, "-o", "Cooking/"+tag+".o")
		oldfileTimings[parameters.file], err = hashTime(t.String())
		checkErr(err)
	}

	hashJSONnew.Body.Entity = append(hashJSONnew.Body.Entity, entity{File: parameters.file, Hash: oldfileTimings[parameters.file]})

	for _, name := range parameters.deps {
		compareAndCompile(name)
	}
}

func linkAll() {
	//compile all the generated .o files under the Cooking directory

	toBeExecuted := "-o " + compilerDetails.name + " " + compilerDetails.includes +
		" " + compilerDetails.otherFlags + " " + "*.o" + " " + compilerDetails.ldFlags

	//exec.Command(compilerDetails.binary, toBeExecuted)
	doNothing(toBeExecuted)
}
