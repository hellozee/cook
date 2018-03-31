package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

//Simple Error Checker
func checkErr(err error) {
	if err != nil {
		panic(err)
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
