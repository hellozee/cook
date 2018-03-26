package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type compiler struct {
	binary     string
	name       string
	start      string
	ldFlags    string
	includes   string
	otherFlags string
}

type params struct {
	file string
	deps []string
}

type entity struct {
	File string `json:"file"`
	Hash string `json:"hash"`
}

type parent struct {
	Body struct {
		Entity []entity `json:"entity"`
	} `json:"body"`
}

//Never Liked Global variables but until I think of a workaround
var compilerDetails compiler
var fileDetails map[string]params
var newfileTimings map[string]string
var oldfileTimings map[string]string
var hashJSONold parent
var hashJSONnew parent
var fileList map[string]string

//Stop Go from throwing warnings if a variable is not used
func doNothing(str string) {
	//Go is badass
}

//Simple Error Checker
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

//Generate the list of files to be compiled
func generateFileList(tag string) {
	details := fileDetails[tag]

	_, err := os.Stat(details.file)
	checkErr(err)

	fileList[tag] = details.file

	if details.deps == nil {
		return
	}

	for _, name := range details.deps {
		generateFileList(name)
	}
}

//Function for executing and debugging exec.Cmd
func checkCommand(cmd *exec.Cmd) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}
}

//Generating hash from timestamp
func hashTime(timeStamp string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(timeStamp), 14)
	return string(bytes), err
}

//Comparing hashes of the current timestamp with the previous one
func checkTimeStamp(timeStamp string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(timeStamp))
	return err == nil
}

func fillCompilerDetails(identifier string, param string) {
	if identifier == "binary" {
		compilerDetails.binary = param
	}
	if identifier == "name" {
		compilerDetails.name = param
	}
	if identifier == "start" {
		compilerDetails.start = param
	}
	if identifier == "ldflags" {
		compilerDetails.ldFlags = param
	}
	if identifier == "includes" {
		compilerDetails.includes = param
	}
	if identifier == "others" {
		compilerDetails.otherFlags = param
	}
}

func fillFileDetails(name string, identifier string, param string) {
	var temp params

	if identifier == "file" {
		temp.file = param
	} else if param != "" {
		temp = fileDetails[name]
	}

	if param == "" {
		return
	}

	if identifier == "deps" {
		paramArray := strings.Split(param, " ")
		temp.deps = paramArray
	}

	fileDetails[name] = temp
}
