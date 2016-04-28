package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/zazab/zhash"
)

func getHash(path string) (zhash.Hash, error) {
	data := map[string]interface{}{}
	_, err := toml.DecodeFile(path, &data)
	if err != nil {
		return zhash.Hash{}, err
	}

	return zhash.HashFromMap(data), nil
}

func GetConfig(
	path string,
) (user, pass string, aliases map[string]string, err error) {
	hash, err := getHash(path)
	if err != nil {
		return "", "", nil, err
	}

	user, err = hash.GetString("user")
	if err != nil {
		return "", "", nil, err
	}

	pass, err = hash.GetString("pass")
	if err != nil {
		return "", "", nil, err
	}

	aliasesData, err := hash.GetMap("aliases")
	if err != nil && !zhash.IsNotFound(err) {
		return "", "", nil, err
	}

	aliases = map[string]string{}
	for alias, hook := range aliasesData {
		aliases[alias] = fmt.Sprintf("%s", hook)
	}

	return user, pass, aliases, nil
}
