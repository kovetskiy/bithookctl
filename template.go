package main

import "bufio"

import "fmt"

import "os"
import "text/template"

func getHookTemplateFunctions(vars map[string]string) template.FuncMap {
	return template.FuncMap{
		"_": func(name string) string {
			value, ok := vars[name]
			if ok {
				return value
			}

			fmt.Printf("Enter template variable '%s' value: ", name)
			value, _ = bufio.NewReader(os.Stdin).ReadString('\n')

			return value
		},
	}
}
