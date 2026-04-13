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
	languageArg := func(lang langs.Lang) string {
		return "--lang=" + string(lang)
	}
	for _, arg := range os.Args[1:] {
		switch arg {
		case languageArg(langs.PT_BR):
			language = langs.PT_BR
		case languageArg(langs.CH):
			language = langs.CH
		case languageArg(langs.JAP):
			language = langs.JAP
		case languageArg(langs.RU):
			language = langs.RU
		default:
			filteredArgs = append(filteredArgs, arg)
		}
	}
	os.Args = filteredArgs

	r := runner.New(language, execCommand)

	r.AutoExecute()
}

func execCommand(cmd string, args ...string) runner.Exec {
	return exec.Command(cmd, args...)
}
