package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kovetskiy/bithooks"
	"github.com/reconquest/executil-go"
	"github.com/reconquest/hierr-go"
)

func handleModeUpdate(
	api *API, hooks *bithooks.Hooks,
	hookName, hookID string,
) error {
	hook, found := hooks.Get(hookName, hookID)
	if !found {
		return errors.New("hook with specified name and id not defined")
	}

	tempDir, err := ioutil.TempDir(os.TempDir(), hookName+"@"+hookID)
	if err != nil {
		return hierr.Errorf(err, "can't create temporary directory")
	}

	argsPath := filepath.Join(tempDir, hookID+"."+hookName+".conf")

	err = ioutil.WriteFile(
		argsPath, []byte(strings.Join(hook.Args, "\n")), 0755,
	)
	if err != nil {
		return hierr.Errorf(err, "can't write hook args")
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, argsPath)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		if executil.IsExitError(err) {
			code := executil.GetExitStatus(err)
			if code != 0 {
				log.Printf(
					"editor exited with non-zero exit code (%d), exiting too",
					code,
				)
				os.Exit(128)
			}
		}

		return hierr.Errorf(err, "exec [%q] failed", cmd.Args)
	}

	argsData, err := ioutil.ReadFile(argsPath)
	if err != nil {
		return hierr.Errorf(err, "can't read hook args file")
	}

	hook.Args = strings.Split(string(argsData), "\n")

	return nil
}
