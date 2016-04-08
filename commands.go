package main

import (
	"strconv"
	"strings"

	"github.com/fsouza/go-dockerclient"
)

type Commands struct {
}

func (r *Commands) Revive() error {
	parrot.Debug("Reviving gadget will reinitialize all the datas.")

	if err := repository.BackupSchema(); err != nil {
		parrot.Warn("Impossible to backup the schema")
		return err
	}

	repository.InitSchema()
	parrot.Println("Gadget revived")
	return nil
}

func (r *Commands) Update() error {
	endpoint := settings.LocalDockerEndpoint()
	client, err := docker.NewClient(endpoint)

	if err != nil {
		return err
	}

	if settings.UseDockerMachine() {
		endpoint = settings.MachineDockerEndpoint()

		parrot.Debug("Configuring client for tls")

		client, _ = docker.NewTLSClient(endpoint,
			settings.MachineDockerCertFile(),
			settings.MachineDockerKeyFile(),
			settings.MachineDockerCAFile())
	}

	imgs, err := client.ListImages(docker.ListImagesOptions{All: false})

	if err != nil {
		return err
	}

	var c = 0

	for _, img := range imgs {
		var id = TruncateID(img.ID)

		parrot.Debug("ID is", id)

		if !repository.Exists(id) {
			repository.Put(img)
			c = c + 1
			parrot.Debug("["+id+"] - ", strings.Join(img.RepoTags, ", "), "added to bucket")
		} else {
			parrot.Debug("["+id+"] - ", strings.Join(img.RepoTags, ", "), " not inserted in bucket because already exists")
		}
	}

	parrot.Println("Added " + strconv.Itoa(c) + " images")

	return nil
}
