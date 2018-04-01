package worker

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	mg "github.com/hellozee/cook/manager"
	ps "github.com/hellozee/cook/parser"
)

//Worker  Data Structure to hold flags necessary for the worker object
type Worker struct {
	QuietFlag   bool
	VerboseFlag bool
}

//CompileFirst  Function for compiling the files for the first time
func (wor *Worker) CompileFirst(par ps.Parser, man mg.Manager) error {
	//Iteratively generate .o files

	for key, value := range man.FileList {
		if wor.QuietFlag == false {
			fmt.Println("Compiling " + value)
		}
		cmd := exec.Command(par.CompilerDetails.Binary, "-c", value,
			"-o", "Cooking/"+key+".o")
		err := checkCommand(cmd, wor)

		if err != nil {
			return err
		}
	}

	return nil
}

/*CompareAndCompile  Function for comparing the hash and the compiling if
the hash did not match */
func (wor *Worker) CompareAndCompile(par ps.Parser, man *mg.Manager) error {
	for key, value := range man.FileList {

		file, err := ioutil.ReadFile(key)

		if err != nil {
			return err
		}

		if !mg.CheckHash(file, man.OldFileTimings[value]) {

			if wor.QuietFlag == false {
				fmt.Println("Compiling " + value)
			}
			cmd := exec.Command(par.CompilerDetails.Binary, "-c", value,
				"-o", "Cooking/"+key+".o")
			err = checkCommand(cmd, wor)

			if err != nil {
				return err
			}

			man.OldFileTimings[value] = mg.HashFile(file)

		}

		man.HashJSONnew.Body.Entity = append(man.HashJSONnew.Body.Entity,
			mg.Entity{File: value, Hash: man.OldFileTimings[value]})
	}

	return nil
}

//Link  Function to link the object files generated
func (wor *Worker) Link(par ps.Parser) error {

	//Compile all the generated .o files under the Cooking directory
	if wor.QuietFlag == false {
		fmt.Println("Linking files..")
	}
	args := []string{par.CompilerDetails.Binary, "-o", par.CompilerDetails.Name,
		par.CompilerDetails.Includes, par.CompilerDetails.OtherFlags,
		"Cooking/*.o", par.CompilerDetails.LdFlags}
	cmd := exec.Command(os.Getenv("SHELL"), "-c", strings.Join(args, " "))
	err := checkCommand(cmd, wor)

	return err
}

//checkCommand  Function to run a command and report any errors
func checkCommand(cmd *exec.Cmd, wor *Worker) error {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if wor.VerboseFlag == true {
		fmt.Println(strings.Join(cmd.Args, " "))
	}
	if err != nil {
		return err
	}

	return nil
}
