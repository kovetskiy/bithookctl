package main

import "github.com/kovetskiy/bithooks"

func handleModeRemove(
	api *API, hooks *bithooks.Hooks,
	hookName, hookID string,
) error {
	hooks.Delete(hookName, hookID)

	return nil
}
