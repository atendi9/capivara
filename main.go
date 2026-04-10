package main

import (
	"os"
	"os/exec"

	"github.com/atendi9/capivara/langs"
	"github.com/atendi9/capivara/runner" 
)

func main() {
	language := langs.EN_US

	var filteredArgs []string
	filteredArgs = append(filteredArgs, os.Args[0])

	for _, arg := range os.Args[1:] {
		if arg == "--lang=portuguese" {
			language = langs.PT_BR
		} else {
			filteredArgs = append(filteredArgs, arg)
		}
	}
	os.Args = filteredArgs

	r := runner.New(language, execCommand)
	
	r.Execute()
}

func execCommand(cmd string, args ...string) runner.Exec {
	return exec.Command(cmd, args...)
}