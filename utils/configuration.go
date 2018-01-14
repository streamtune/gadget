package utils

import (
	"encoding/json"
	"path/filepath"

	"github.com/gi4nks/quant"
)

type Configuration struct {
	parrot *quant.Parrot

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

func NewConfiguration(p quant.Parrot) *Configuration {
	var c = Configuration{}
	c.parrot = &p

	c.RepositoryDirectory = ConstRepositoryDirectory
	c.RepositoryFile = ConstRepositoryFile
	c.DebugMode = ConstDebugMode
	c.RestServerMode = ConstRestServerMode
	c.RestPort = ConstRestPort
	c.LocalDockerEndpoint = ConstLocalDockerEndpoint
	c.MachineDockerEndpoint = ConstMachineDockerEndpoint
	c.UseDockerMachine = ConstUseDockerMachine
	c.MachineDockerCertFile = ConstMachineDockerCertFile
	c.MachineDockerKeyFile = ConstMachineDockerKeyFile
	c.MachineDockerCAFile = ConstMachineDockerCAFile

	return &c
}

func (c Configuration) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		c.parrot.Error("Warning", err)
		return "{}"
	}
	return string(b)
}

func (c Configuration) RepositoryFullName() string {

	return c.RepositoryDirectory + string(filepath.Separator) + c.RepositoryFile
}
