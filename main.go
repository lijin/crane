package main

import (
	"errors"
	"fmt"
	"github.com/michaelsauter/crane/print"
	"os"
	"os/exec"
	"strings"
	"path/filepath"
)

type StatusError struct {
	error  error
	status int
}

func main() {
	// On panic, recover the error, display it and return the given status code if any
	defer func() {
		var statusError StatusError

		switch err := recover().(type) {
		case StatusError:
			statusError = err
		case error:
			statusError = StatusError{err, 1}
		case string:
			statusError = StatusError{errors.New(err), 1}
		default:
			statusError = StatusError{}
		}

		if statusError.error != nil {
			print.Error("ERROR: %s\n", statusError.error)
		}
		os.Exit(statusError.status)
	}()

	handleCmd()
}

// from https://stackoverflow.com/questions/20437336/how-to-execute-system-command-in-golang-with-unknown-arguments
var (
    output_path = filepath.Join("./output")
    bash_script = filepath.Join( "_script.sh" )
)
func checkError( e error){
    if e != nil {
        panic(e)
    }
}
func exe_cmd(cmds []string) {
    os.RemoveAll(output_path)
    err := os.MkdirAll( output_path, os.ModePerm|os.ModeDir )
    checkError(err)
    file, err := os.Create( filepath.Join(output_path, bash_script))
    checkError(err)
    defer file.Close()
    file.WriteString("#!/bin/sh\n")
    file.WriteString( strings.Join(cmds, " "))
    err = os.Chdir(output_path)
    checkError(err)
    out, err := exec.Command("sh", bash_script).Output()
    checkError(err)
    fmt.Println(string(out))
}

func executeCommand(name string, args []string) {
	if isVerbose() {
		fmt.Printf("\n--> %s %s\n", name, strings.Join(args, " "))
	}
	exe_cmd(append([]string{name}, args...))
}

func commandOutput(name string, args []string) (string, error) {
	out, err := exec.Command(name, args...).CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// from https://gist.github.com/dagoof/1477401
func pipedCommandOutput(pipedCommandArgs ...[]string) ([]byte, error) {
	var commands []exec.Cmd
	for _, commandArgs := range pipedCommandArgs {
		cmd := exec.Command(commandArgs[0], commandArgs[1:]...)
		commands = append(commands, *cmd)
	}
	for i, command := range commands[:len(commands)-1] {
		out, err := command.StdoutPipe()
		if err != nil {
			return nil, err
		}
		command.Start()
		commands[i+1].Stdin = out
	}
	final, err := commands[len(commands)-1].Output()
	if err != nil {
		return nil, err
	}
	return final, nil
}
