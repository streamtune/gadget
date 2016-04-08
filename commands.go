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
			// load image details

			// skipping dangling images
			if img.RepoTags[0] != "<none>:<none>" {

				imgDetails, err := client.InspectImage(img.RepoTags[0])

				if err != nil {
					return err
				}

				repository.Put(img, *imgDetails)
			} else {
				repository.Put(img, docker.Image{})
			}
			c = c + 1
			parrot.Debug("["+id+"] - ", strings.Join(img.RepoTags, ", "), "added to bucket")
		} else {
			parrot.Debug("["+id+"] - ", strings.Join(img.RepoTags, ", "), " not inserted in bucket because already exists")
		}
	}

	parrot.Println("Added " + strconv.Itoa(c) + " images")

	return nil
}

func (r *Commands) List() {

	var images = repository.GetAll()

	var header = Image{}
	var iis = [][]string{}

	for _, i := range images {
		for _, r := range i.RowsForList() {
			iis = append(iis, r)
		}
	}

	parrot.TablePrint(header.HeaderForList(), iis)

}

func (r *Commands) Labels() {
	var images = repository.GetAll()

	var header = Image{}
	var iis = [][]string{}

	for _, i := range images {
		for _, r := range i.RowsForLabel() {
			iis = append(iis, r)
		}
	}

	parrot.TablePrint(header.HeaderForLabel(), iis)
}

func (r *Commands) LabelsById(id string) {
	var image = repository.Get(id)

	var header = Image{}
	var iis = [][]string{}

	for _, r := range image.RowsForLabel() {
		iis = append(iis, r)
	}

	parrot.TablePrint(header.HeaderForLabel(), iis)
}

func (r *Commands) LabelsByTag(id string) {
	var image = repository.FindByTag(id)

	var header = Image{}
	var iis = [][]string{}

	for _, r := range image.RowsForLabel() {
		iis = append(iis, r)
	}

	parrot.TablePrint(header.HeaderForLabel(), iis)
}

func (r *Commands) Info() {

	var images = repository.GetAll()

	var header = Image{}
	var iis = [][]string{}

	for _, i := range images {
		for _, r := range i.RowsForInfo() {
			iis = append(iis, r)
		}
	}

	parrot.TablePrint(header.HeaderForInfo(), iis)
}

func (r *Commands) InfoById(id string) {
	var image = repository.Get(id)

	var header = Image{}
	var iis = [][]string{}

	for _, r := range image.RowsForInfo() {
		iis = append(iis, r)
	}

	parrot.TablePrint(header.HeaderForInfo(), iis)
}

func (r *Commands) InfoByTag(id string) {
	var image = repository.FindByTag(id)

	var header = Image{}
	var iis = [][]string{}

	for _, r := range image.RowsForInfo() {
		iis = append(iis, r)
	}

	parrot.TablePrint(header.HeaderForInfo(), iis)
}

func (r *Commands) Volumes() {
	var images = repository.GetAll()

	var header = Image{}
	var iis = [][]string{}

	for _, i := range images {
		for _, r := range i.RowsForVolume() {
			iis = append(iis, r)
		}
	}

	parrot.TablePrint(header.HeaderForVolume(), iis)
}

func (r *Commands) VolumesById(id string) {
	var image = repository.Get(id)

	var header = Image{}
	var iis = [][]string{}

	for _, r := range image.RowsForVolume() {
		iis = append(iis, r)
	}

	parrot.TablePrint(header.HeaderForVolume(), iis)
}

func (r *Commands) VolumesByTag(id string) {
	var image = repository.FindByTag(id)

	var header = Image{}
	var iis = [][]string{}

	for _, r := range image.RowsForVolume() {
		iis = append(iis, r)
	}

	parrot.TablePrint(header.HeaderForVolume(), iis)
}
