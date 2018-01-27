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

type entity struct {
	File string `json:"file"`
	Hash string `json:"hash"`
}

type parent struct {
	Body struct {
		Entity []entity `json:"entity"`
	} `json:"body"`
}

var compilerDetails compiler
var fileDetails map[string]params
var newfileTimings map[string]string
var oldfileTimings map[string]string
var hashJSONold parent
var hashJSONnew parent

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
