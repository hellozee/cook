package worker

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	lg "github.com/hellozee/cook/logger"
	mg "github.com/hellozee/cook/manager"
	ps "github.com/hellozee/cook/parser"
)

//Worker  Data Structure to hold flags necessary for the worker object
type Worker struct {
	Logger *lg.Logger
}

//CompileFirst  Function for compiling the files for the first time
func (wor *Worker) CompileFirst(par ps.Parser, man mg.Manager) error {
	//Iteratively generate .o files

	for key, value := range man.FileList {
		wor.Logger.ReportSuccess("Compiling " + value)
		cmd := exec.Command(par.CompilerDetails.Binary, "-c", value,
			"-o", "Cooking/"+key+".o")
		err := checkCommand(cmd, wor)

		if err != nil {
			wor.Logger.ReportError(err.Error())
			return err
		}
	}
	wor.Logger.ReportSuccess("Successfully Compiled all the files")
	return nil
}

/*CompareAndCompile  Function for comparing the hash and the compiling if
  the hash did not match */
func (wor *Worker) CompareAndCompile(par ps.Parser, man *mg.Manager) error {
	for key, value := range man.FileList {

		file, err := ioutil.ReadFile(key)

		if err != nil {
			wor.Logger.ReportError(err.Error())
			return err
		}

		if !mg.CheckHash(file, man.OldFileTimings[value]) {

			wor.Logger.ReportSuccess("Compiling " + value)
			cmd := exec.Command(par.CompilerDetails.Binary, "-c", value,
				"-o", "Cooking/"+key+".o")
			err = checkCommand(cmd, wor)

			if err != nil {
				wor.Logger.ReportError(err.Error())
				return err
			}

			man.OldFileTimings[value] = mg.HashFile(file)

		}

		man.HashJSONnew.Body.Entity = append(man.HashJSONnew.Body.Entity,
			mg.Entity{File: value, Hash: man.OldFileTimings[value]})
	}
	wor.Logger.ReportSuccess("Successfully Compiled all the files")
	return nil
}

//Link  Function to link the object files generated
func (wor *Worker) Link(par ps.Parser) error {

	//Compile all the generated .o files under the Cooking directory
	wor.Logger.ReportSuccess("Linking files")
	args := []string{par.CompilerDetails.Binary, "-o", par.CompilerDetails.Name,
		par.CompilerDetails.Includes, par.CompilerDetails.OtherFlags,
		"Cooking/*.o", par.CompilerDetails.LdFlags}
	cmd := exec.Command(os.Getenv("SHELL"), "-c", strings.Join(args, " "))
	err := checkCommand(cmd, wor)
	wor.Logger.ReportSuccess("Successfully Linked files")
	return err
}

//checkCommand  Function to run a command and report any errors
func checkCommand(cmd *exec.Cmd, wor *Worker) error {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	wor.Logger.ReportWarning(out.String())
	wor.Logger.ReportError(stderr.String())
	if err != nil {
		wor.Logger.ReportError(err.Error())
		return err
	}
	wor.Logger.ReportSuccess("Ran Successfully " + strings.Join(cmd.Args, " "))
	return nil
}

//NewWorker  Function to create a new worker
func NewWorker(log *lg.Logger) Worker {
	wor := Worker{
		Logger: log,
	}

	return wor
}
