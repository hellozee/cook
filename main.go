package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func init() {
	fileDetails = make(map[string]params)
	newfileTimings = make(map[string]string)
	oldfileTimings = make(map[string]string)
}

func generateList(tag string) {

	parameters := fileDetails[tag]

	file, err := os.Stat(parameters.file)
	checkErr(err)
	t := file.ModTime()
	hash, err := hashTime(t.String())
	checkErr(err)
	newfileTimings[parameters.file] = hash
	hashJSONnew.Body.Entity = append(hashJSONnew.Body.Entity, entity{File: parameters.file, Hash: hash})
	for _, name := range parameters.deps {
		generateList(name)
	}

}

func main() {

	temp, err := ioutil.ReadFile("Recipe")
	Recipe := string(temp)
	checkErr(err)

	parse(Recipe)

	if _, err := os.Stat("Cooking/example.json"); err == nil {
		fmt.Println("It exists")
		xmlFile, err := os.Open("Cooking/example.json")
		checkErr(err)
		defer xmlFile.Close()

		bytes, _ := ioutil.ReadAll(xmlFile)
		var parsed parent
		err = json.Unmarshal(bytes, &parsed)
		fmt.Println(parsed)

		generateList(compilerDetails.start)

		//compare();

	} else {
		_ = os.Mkdir("Cooking", 0777)
		generateList(compilerDetails.start)

		//compileFirst();

	}

}
