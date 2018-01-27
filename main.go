package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func init() {
	fileDetails = make(map[string]params)
	newfileTimings = make(map[string]string)
	oldfileTimings = make(map[string]string)
}

func main() {

	temp, err := ioutil.ReadFile("Recipe")
	Recipe := string(temp)
	checkErr(err)

	parse(Recipe)

	var jsonData []byte

	if _, err := os.Stat("Cooking/details.json"); err == nil {

		jsonFile, err := os.Open("Cooking/details.json")
		checkErr(err)
		defer jsonFile.Close()
		bytes, _ := ioutil.ReadAll(jsonFile)
		err = json.Unmarshal(bytes, &hashJSONold)
		checkErr(err)

		structToMap(hashJSONold.Body.Entity)

		compareAndCompile(compilerDetails.start)

	} else {
		_ = os.Mkdir("Cooking", 0777)
		generateList(compilerDetails.start)
		compileFirst(compilerDetails.start)
	}

	jsonData, err = json.MarshalIndent(hashJSONnew, "", " ")
	checkErr(err)

	file, err := os.OpenFile("Cooking/details.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0777)
	checkErr(err)
	defer file.Close()

	_, err = file.Write(jsonData)
	checkErr(err)

}
