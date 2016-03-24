package main

import (
	"encoding/json"
	"io/ioutil"
)

type Configuration struct {
	RepositoryDirectory string
	RepositoryFile      string
	DebugMode           bool
	RestServerMode      string
	RestPort            int
}

type Settings struct {
	configs Configuration
}

func (sts *Settings) LoadSettings() {
	folder := executableFolder()

	file, err := ioutil.ReadFile(folder + "/conf.json")

	if err != nil {
		sts.configs = Configuration{}
		sts.configs.RepositoryDirectory = folder + "/" + ConstRepositoryDirectory
		sts.configs.RepositoryFile = ConstRepositoryFile
		sts.configs.DebugMode = ConstDebugMode
		sts.configs.RestServerMode = ConstRestServerMode
		sts.configs.RestPort = ConstRestPort

	} else {
		json.Unmarshal(file, &sts.configs)

		parrot.Debug("folder: " + folder)
		parrot.Debug("file: " + asJson(sts.configs))

	}
}

func (sts Settings) RepositoryDirectory() string {
	return sts.configs.RepositoryDirectory
}

func (sts Settings) RepositoryFile() string {
	return sts.configs.RepositoryFile
}

func (sts Settings) DebugMode() bool {
	return sts.configs.DebugMode
}

func (sts Settings) RestServerMode() string {
	return sts.configs.RestServerMode
}

func (sts Settings) RestPort() int {
	return sts.configs.RestPort
}

func (sts Settings) String() string {
	b, err := json.Marshal(sts.configs)
	if err != nil {
		parrot.Error("Warning", err)
		return "{}"
	}
	return string(b)
}
