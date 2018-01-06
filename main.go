package main

import "os"
import "io/ioutil"

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

	dir, err := os.Getwd()
	checkErr(err)
	doNothing(dir)

	file, err := os.Stat("Recipe")
	checkErr(err)

	t := file.ModTime()

	doNothing(t.String())
}
