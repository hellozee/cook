package main

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"os"
	"os/exec"
	"strings"

	ps "github.com/hellozee/cook/parser"
)

type entity struct {
	File string `json:"file"`
	Hash uint32 `json:"hash"`
}

type parent struct {
	Body struct {
		Entity []entity `json:"entity"`
	} `json:"body"`
}

//Never Liked Global variables but until I think of a workaround
var newfileTimings map[string]uint32
var oldfileTimings map[string]uint32
var hashJSONold parent
var hashJSONnew parent
var fileList map[string]string

//Simple Error Checker
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

//Generate the list of files to be compiled
func generateFileList(par ps.Parser, tag string) {
	details := par.FileDetails[tag]

	_, err := os.Stat(details.File)
	checkErr(err)

	fileList[tag] = details.File

	if details.Deps == nil {
		return
	}

	for _, name := range details.Deps {
		generateFileList(par, name)
	}
}

//Function for executing and debugging exec.Cmd
func checkCommand(cmd *exec.Cmd) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if *verboseFlag == true {
		fmt.Println(strings.Join(cmd.Args, " "))
	}
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}
}

//Generating hash from timestamp
func hashFile(file []byte) uint32 {
	hash := crc32.ChecksumIEEE(file)
	return hash
}

//Comparing hashes of the current timestamp with the previous one
func checkHash(file []byte, hash uint32) bool {
	generatedHash := crc32.ChecksumIEEE(file)
	return generatedHash == hash
}
