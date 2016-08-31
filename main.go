package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kovetskiy/bithooks"
	"github.com/kovetskiy/godocs"
	"github.com/seletskiy/hierr"
)

var (
	version = "1.0"
	usage   = os.ExpandEnv(
		`bithookctl ` + version + `

Manage your hooks in Bitbucket (Atlassian Stash) repository.

You should create configuration file at ~/.config/bithookctl.conf
with following syntax:

  user = "username"
  pass = "password"

If you will work with multiple hooks you can add aliases section with following
syntax:

  [aliases]
   pre = "com.ngs.stash.externalhooks.external-hooks:external-pre-receive-hook"
   post = "com.ngs.stash.externalhooks.external-hooks:external-post-receive-hook"

Usage:
    bithookctl [options] -L
    bithookctl [options] -A <hook> <id> [(-v <var_name> <var_value>)]...
    bithookctl [options] -U <hook> <id>
    bithookctl [options] -R <hook> <id>
    bithookctl -h | --help

Options:
    -L --list       List installed hooks.
    -A --add        Add hook <hook> with <id>, compile <templates>/<hook>
                     template and use as <hook> args.
    -U --update     Update hook <hook> with <id> args.
    -R --remove     Remove hook <id> with <id>.
    -u <url>        Specify repository URL.
                     By default, it will be read from 'git remote origin'
                     output.
    -c <config>     Specify configuration file with user credentials.
                     [default: $HOME/.config/bithookctl.conf]
    -t <templates>  Specify directory with templates.
                     [default: /var/lib/bithookctl/templates/]
    -p <key>        Specify Bitbucket hook key or an alias for hook key in
                     [hooks] section from configuration file.
                     [default: pre]
    -v              Set template variable <var_name> value to <var_value>.
    -h --help       Show this screen.
    --version       Show version.
`)
)

func main() {
	args := godocs.MustParse(usage, version, godocs.UsePager)

	log.SetFlags(0)

	var (
		repoURL, _    = args["-u"].(string)
		configPath    = args["-c"].(string)
		templatesRoot = args["-t"].(string)
		hookKey       = args["-p"].(string)

		hookName, _ = args["<hook>"].(string)
		hookID, _   = args["<id>"].(string)

		varsNames, _  = args["<var_name>"].([]string)
		varsValues, _ = args["<var_value>"].([]string)

		mode = getMode(args)
	)

	user, pass, aliases, err := GetConfig(configPath)
	if err != nil {
		hierr.Fatalf(err, "can't load configuration")
	}

	repo, err := GetRepo(repoURL)
	if err != nil {
		hierr.Fatalf(err, "can't get repository")
	}

	if value, ok := aliases[hookKey]; ok {
		hookKey = value
	}

	api := NewAPI(repo, user, pass, hookKey)

	settings, err := api.GetHookSettings()
	if err != nil {
		hierr.Fatalf(err, "can't get hook settings")
	}

	if mode == "list" {
		fmt.Print("Hooks:\n" + settings.Params)
		os.Exit(0)
	}

	if settings.Exe != "bithooker" {
		settings.Exe = "bithooker"
		settings.Params = ""
	}

	hooks, err := bithooks.Decode(settings.Params)
	if err != nil {
		hierr.Fatalf(
			hierr.Errorf(settings.Params, err.Error()),
			"can't decode hook params",
		)
	}

	switch mode {
	case "add":
		err = handleModeAdd(
			api, &hooks, hookName, hookID, templatesRoot,
			generateTemplateVars(varsNames, varsValues),
		)

	case "update":
		err = handleModeUpdate(api, &hooks, hookName, hookID)

	case "remove":
		err = handleModeRemove(api, &hooks, hookName, hookID)

	default:
		log.Fatalln("unexpected mode")
	}

	if err != nil {
		log.Fatalln(err.Error())
	}

	settings.Params = bithooks.Encode(hooks)
	settings.SafePath = true

	err = api.SetHookSettings(settings)
	if err != nil {
		hierr.Fatalf(err, "can't save hook settings")
	}

	err = api.EnableHook()
	if err != nil {
		hierr.Fatalf(err, "can't enable hook")
	}

	fmt.Println("hook settings saved")
}

func generateTemplateVars(names, values []string) map[string]string {
	vars := map[string]string{}
	for index, name := range names {
		vars[name] = values[index]
	}

	return vars
}

func getMode(args map[string]interface{}) string {
	var (
		modes = []string{
			"list", "add", "update", "remove",
		}
	)

	for _, mode := range modes {
		modeValue, _ := args["--"+mode].(bool)
		if modeValue {
			return mode
		}
	}

	return ""
}
