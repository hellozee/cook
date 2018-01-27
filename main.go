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

func main() {

	temp, err := ioutil.ReadFile("Recipe")
	Recipe := string(temp)
	checkErr(err)

	parse(Recipe)
	fmt.Println(fileDetails)

	if _, err := os.Stat("Cooking/details.json"); err == nil {
		fmt.Println("It exists")
		xmlFile, err := os.Open("Cooking/details.json")
		checkErr(err)
		defer xmlFile.Close()
		bytes, _ := ioutil.ReadAll(xmlFile)
		var parsed parent
		err = json.Unmarshal(bytes, &parsed)

		fmt.Println(parsed)

		//compareAndCompile();

		generateList(compilerDetails.start)

	} else {
		_ = os.Mkdir("Cooking", 0777)
		generateList(compilerDetails.start)
		//compileFirst();
	}

	jsonData, err := json.Marshal(hashJSONnew)

	file, err := os.OpenFile("Cooking/details.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0777)
	checkErr(err)
	defer file.Close()

	_, err = file.Write(jsonData)
	checkErr(err)
}
