package main

import (
	"strings"
)

type HookSettings struct {
	Exe      string `json:"exe"`
	SafePath bool   `json:"safe_path"`
	Params   string `json:"params"`
}

type ResponseHook struct {
	Detauls struct {
		Version string `json:"version"`
	} `json:"detauls"`
	Configured bool `json:"configured"`
	Enabled    bool `json:"enabled"`
}

type ResponseError struct {
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func (err ResponseError) String() string {
	var messages []string

	for _, nestedError := range err.Errors {
		messages = append(messages, nestedError.Message)
	}

	return strings.Join(messages, "\n")
}
