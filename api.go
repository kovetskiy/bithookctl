package main

import "github.com/bndr/gopencils"

type API struct {
	*gopencils.Resource
	*Repo
	hookKey string
}

func NewAPI(repo *Repo, user, pass, hookKey string) *API {
	api := &API{
		Resource: gopencils.Api(
			"http://"+repo.host+"/rest/api/1.0"+
				"/"+repo.projectType+"/"+repo.project+"/repos/"+repo.name,
			&gopencils.BasicAuth{user, pass},
		),
		Repo:    repo,
		hookKey: hookKey,
	}

	return api
}

func (api *API) GetHookSettings() (*HookSettings, error) {
	response, err := request(
		"GET",
		api.Res("settings").
			Res("hooks").Res(api.hookKey).
			Res("settings", &HookSettings{}),
	)
	if err != nil {
		return nil, err
	}

	return response.(*HookSettings), nil
}

func (api *API) IsHookEnabled() (bool, error) {
	rawResponse, err := request(
		"GET",
		api.Res("settings").
			Res("hooks").Res(api.hookKey, &ResponseHook{}),
	)
	if err != nil {
		return false, err
	}

	response := rawResponse.(*ResponseHook)
	return response.Enabled, nil
}

func (api *API) EnableHook() error {
	_, err := request(
		"PUT",
		api.Res("settings").
			Res("hooks").Res(api.hookKey).
			Res("enabled", &ResponseHook{}),
	)

	return err
}

func (api *API) SetHookSettings(settings *HookSettings) error {
	_, err := request(
		"PUT",
		api.Res("settings").
			Res("hooks").Res(api.hookKey).
			Res("settings", &HookSettings{}),
		settings,
	)

	return err
}
