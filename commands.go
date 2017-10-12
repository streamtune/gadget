package main

import (
	"strings"

	"github.com/fsouza/go-dockerclient"
	"gopkg.in/cheggaaa/pb.v1"
)

type Commands struct {
}

func (r *Commands) Debug(dbg bool) {
	parrot.Debug("Changing debug mode.")

	settings.SetDebugMode(dbg)

	parrot.Println("Switched debug mode to", dbg)
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

	count := len(imgs)
	bar := pb.StartNew(count)

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

		bar.Increment()
	}

	bar.FinishPrint("Done.")

	return nil
}

func (r *Commands) List() {

	images, err := repository.GetAll()

	if err != nil {
		parrot.Error("Error", err)
		return
	}

	var header = Image{}
	var iis = [][]string{}

	for _, i := range images {
		for _, r := range i.RowsForList() {
			iis = append(iis, r)
		}
	}

	parrot.Tablify(header.HeaderForList(), iis)

}

func (r *Commands) ListByNumber(co int) {
	var header = Image{}

	if co <= 0 {
		parrot.Tablify(header.HeaderForList(), nil)
		return
	}

	images, err := repository.GetLimit(co)

	if err != nil {
		parrot.Error("Error", err)
		return
	}

	var iis = [][]string{}

	for _, i := range images {
		for _, r := range i.RowsForList() {
			iis = append(iis, r)
		}
	}

	parrot.Tablify(header.HeaderForList(), iis)

}

func (r *Commands) ListByName(name string) {
	imagesMap := make(map[string]bool)
	images, err := repository.GetAll()

	if err != nil {
		parrot.Error("Error", err)
		return
	}

	var header = Image{}
	var iis = [][]string{}

	for _, i := range images {
		for _, t := range i.Tags {
			if strings.Contains(t.ID, name) && !imagesMap[i.ID] {
				imagesMap[i.ID] = true

				for _, r := range i.RowsForList() {
					iis = append(iis, r)
				}
			}
		}
	}

	parrot.Tablify(header.HeaderForList(), iis)

}

func (r *Commands) Labels() {
	images, err := repository.GetAll()

	if err != nil {
		parrot.Error("Error", err)
		return
	}

	var header = Image{}
	var iis = [][]string{}

	for _, i := range images {
		for _, r := range i.RowsForLabel() {
			iis = append(iis, r)
		}
	}

	parrot.Tablify(header.HeaderForLabel(), iis)
}

func (r *Commands) LabelsById(id string) {
	image, err := repository.Get(id)

	if err != nil {
		parrot.Error("Error", err)
		return
	}

	var header = Image{}
	var iis = [][]string{}

	for _, r := range image.RowsForLabel() {
		iis = append(iis, r)
	}

	parrot.Tablify(header.HeaderForLabel(), iis)
}

func (r *Commands) LabelsByTag(id string) {
	images, err := repository.FindByTag(id)

	if err != nil {
		parrot.Error("Error", err)
		return
	}

	var header = Image{}
	var iis = [][]string{}

	for _, i := range images {
		for _, r := range i.RowsForLabel() {
			iis = append(iis, r)
		}
	}

	parrot.Tablify(header.HeaderForLabel(), iis)
}

func (r *Commands) Info() {

	images, err := repository.GetAll()

	if err != nil {
		parrot.Error("Error", err)
		return
	}

	var header = Image{}
	var iis = [][]string{}

	for _, i := range images {
		for _, r := range i.RowsForInfo() {
			iis = append(iis, r)
		}
	}

	parrot.Tablify(header.HeaderForInfo(), iis)
}

func (r *Commands) InfoById(id string) {
	image, err := repository.Get(id)

	if err != nil {
		parrot.Error("Error", err)
		return
	}

	var header = Image{}
	var iis = [][]string{}

	for _, r := range image.RowsForInfo() {
		iis = append(iis, r)
	}

	parrot.Tablify(header.HeaderForInfo(), iis)
}

func (r *Commands) InfoByTag(id string) {
	images, err := repository.FindByTag(id)

	if err != nil {
		parrot.Error("Error", err)
		return
	}

	var header = Image{}
	var iis = [][]string{}

	for _, i := range images {
		for _, r := range i.RowsForInfo() {
			iis = append(iis, r)
		}
	}

	parrot.Tablify(header.HeaderForInfo(), iis)
}

func (r *Commands) Volumes() {
	images, err := repository.GetAll()

	if err != nil {
		parrot.Error("Error", err)
		return
	}

	var header = Image{}
	var iis = [][]string{}

	for _, i := range images {
		for _, r := range i.RowsForVolume() {
			iis = append(iis, r)
		}
	}

	parrot.Tablify(header.HeaderForVolume(), iis)
}

func (r *Commands) VolumesById(id string) {
	image, err := repository.Get(id)

	if err != nil {
		parrot.Error("Error", err)
		return
	}

	var header = Image{}
	var iis = [][]string{}

	for _, r := range image.RowsForVolume() {
		iis = append(iis, r)
	}

	parrot.Tablify(header.HeaderForVolume(), iis)
}

func (r *Commands) VolumesByTag(id string) {
	images, err := repository.FindByTag(id)

	if err != nil {
		parrot.Error("Error", err)
		return
	}

	var header = Image{}
	var iis = [][]string{}

	for _, i := range images {
		for _, r := range i.RowsForVolume() {
			iis = append(iis, r)
		}
	}
	parrot.Tablify(header.HeaderForVolume(), iis)
}
