package main

import (
	"encoding/json"
	"io/ioutil"
)

type Configuration struct {
	RepositoryDirectory   string
	RepositoryFile        string
	DebugMode             bool
	RestServerMode        string
	RestPort              int
	LocalDockerEndpoint   string
	MachineDockerEndpoint string
	UseDockerMachine      bool
	MachineDockerCertFile string
	MachineDockerKeyFile  string
	MachineDockerCAFile   string
}

type Settings struct {
	configs Configuration
}

func (sts *Settings) LoadSettings() {
	folder, err := pathUtils.ExecutableFolder()

	if err != nil {
		parrot.Error("Executable forlder error", err)
	}

	file, err := ioutil.ReadFile(folder + "/conf.json")

	if err != nil {
		sts.configs = Configuration{}
		sts.configs.RepositoryDirectory = folder + "/" + ConstRepositoryDirectory
		sts.configs.RepositoryFile = ConstRepositoryFile
		sts.configs.DebugMode = ConstDebugMode
		sts.configs.RestServerMode = ConstRestServerMode
		sts.configs.RestPort = ConstRestPort
		sts.configs.LocalDockerEndpoint = ConstLocalDockerEndpoint
		sts.configs.MachineDockerEndpoint = ConstMachineDockerEndpoint
		sts.configs.UseDockerMachine = ConstUseDockerMachine
		sts.configs.MachineDockerCertFile = ConstMachineDockerCertFile
		sts.configs.MachineDockerKeyFile = ConstMachineDockerKeyFile
		sts.configs.MachineDockerCAFile = ConstMachineDockerCAFile

	} else {
		json.Unmarshal(file, &sts.configs)

		parrot.Debug("folder: " + folder)
		parrot.Debug("file: " + AsJson(sts.configs))

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

func (sts Settings) UseDockerMachine() bool {
	return sts.configs.UseDockerMachine
}

func (sts Settings) LocalDockerEndpoint() string {
	return sts.configs.LocalDockerEndpoint
}

func (sts Settings) MachineDockerEndpoint() string {
	return sts.configs.MachineDockerEndpoint
}

func (sts Settings) MachineDockerCertFile() string {
	return sts.configs.MachineDockerCertFile
}

func (sts Settings) MachineDockerKeyFile() string {
	return sts.configs.MachineDockerKeyFile
}

func (sts Settings) MachineDockerCAFile() string {
	return sts.configs.MachineDockerCAFile
}

func (sts Settings) String() string {
	b, err := json.Marshal(sts.configs)
	if err != nil {
		parrot.Error("Warning", err)
		return "{}"
	}
	return string(b)
}
