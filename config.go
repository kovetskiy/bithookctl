package main

import (
	"github.com/BurntSushi/toml"
	"github.com/seletskiy/hierr"
	"github.com/zazab/zhash"
)

func getConfig(path string) (zhash.Hash, error) {
	data := map[string]interface{}{}
	_, err := toml.DecodeFile(path, &data)
	if err != nil {
		return zhash.Hash{}, err
	}

	return zhash.HashFromMap(data), nil
}

func GetUserData(configPath string) (user, pass string, err error) {
	config, err := getConfig(configPath)
	if err != nil {
		return "", "", err
	}

	user, err = config.GetString("user")
	if err != nil {
		return "", "", hierr.Errorf(err, "can't get 'user' field from config")
	}

	pass, err = config.GetString("pass")
	if err != nil {
		return "", "", hierr.Errorf(err, "can't get 'pass' field from config")
	}

	return user, pass, nil
}
