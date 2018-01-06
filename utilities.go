package main

import (
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

var compilerDetails compiler
var fileDetails map[string]params
var newfileTimings map[string]string
var oldfileTimings map[string]string

func doNothing(str string) {
	//Go is badass
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func hashTime(timeStamp string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(timeStamp), 14)
	return string(bytes), err
}

func checkTimeStamp(timeStamp, hash string) bool {
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

	if identifier == "deps" {
		paramArray := strings.Split(param, " ")
		temp.deps = paramArray
	}

	fileDetails[name] = temp
}
