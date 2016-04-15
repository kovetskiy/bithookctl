package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	libtemplate "text/template"

	"github.com/kovetskiy/bithooks"
	"github.com/seletskiy/hierr"
)

func handleModeAdd(
	api *API, hooks *bithooks.Hooks,
	hookName, hookID string,
	templatesDir string, vars map[string]string,
) error {
	_, found := hooks.Get(hookName, hookID)
	if found {
		return fmt.Errorf("hook with specified name and id already defined")
	}

	var (
		hookArgs     = bytes.NewBuffer(nil)
		templateName = hookName + ".template"
		templatePath = filepath.Join(templatesDir, templateName)
	)

	template, err := libtemplate.New(templateName).
		Option("missingkey=error").
		Funcs(getHookTemplateFunctions(vars)).
		ParseFiles(templatePath)
	if err != nil {
		return hierr.Errorf(
			hierr.Errorf(err, templatePath),
			"can't parse template",
		)
	}

	err = template.Execute(hookArgs, vars)
	if err != nil {
		return hierr.Errorf(
			hierr.Errorf(err, templatePath),
			"can't compile template",
		)
	}

	hook := &bithooks.Hook{
		Name: hookName,
		ID:   hookID,
		Args: strings.Split(hookArgs.String(), "\n"),
	}

	err = hooks.Append(hook)
	if err != nil {
		return hierr.Errorf("can't add hook %s@%s", hookName, hookID)
	}

	return nil
}
