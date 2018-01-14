package repos

import (
	"strings"

	"github.com/fsouza/go-dockerclient"
	"gopkg.in/cheggaaa/pb.v1"

	"github.com/gi4nks/quant"

	models "github.com/streamtune/gadget/models"
	utils "github.com/streamtune/gadget/utils"
)

type Commands struct {
	parrot *quant.Parrot

	configuration *utils.Configuration
	repository    *Repository
	utils         *utils.Utilities
}

func NewCommands(p quant.Parrot, c utils.Configuration, r Repository, u utils.Utilities) *Commands {
	return &Commands{parrot: &p, configuration: &c, repository: &r, utils: &u}
}

func (r *Commands) Debug(dbg bool) {
	r.parrot.Debug("Changing debug mode.")

	r.configuration.DebugMode = dbg

	r.parrot.Println("Switched debug mode to", dbg)
}

func (r *Commands) Revive() error {
	r.parrot.Debug("Reviving gadget will reinitialize all the datas.")

	if err := r.repository.BackupSchema(); err != nil {
		r.parrot.Warn("Impossible to backup the schema")
		return err
	}

	r.repository.InitSchema()
	r.parrot.Println("Gadget revived")
	return nil
}

func (r *Commands) Update() error {
	endpoint := r.configuration.LocalDockerEndpoint
	client, err := docker.NewClient(endpoint)

	if err != nil {
		return err
	}

	if r.configuration.UseDockerMachine {
		endpoint = r.configuration.MachineDockerEndpoint

		r.parrot.Debug("Configuring client for tls")

		client, _ = docker.NewTLSClient(endpoint,
			r.configuration.MachineDockerCertFile,
			r.configuration.MachineDockerKeyFile,
			r.configuration.MachineDockerCAFile)
	}

	imgs, err := client.ListImages(docker.ListImagesOptions{All: false})

	if err != nil {
		return err
	}

	var c = 0

	count := len(imgs)
	bar := pb.StartNew(count)

	for _, img := range imgs {
		var id = r.utils.TruncateID(img.ID)

		r.parrot.Debug("ID is", id)

		if !r.repository.Exists(id) {
			// load image details

			// skipping dangling images
			if img.RepoTags[0] != "<none>:<none>" {

				imgDetails, err := client.InspectImage(img.RepoTags[0])

				if err != nil {
					return err
				}

				r.repository.Put(img, *imgDetails)
			} else {
				r.repository.Put(img, docker.Image{})
			}
			c = c + 1
			r.parrot.Debug("["+id+"] - ", strings.Join(img.RepoTags, ", "), "added to bucket")
		} else {
			r.parrot.Debug("["+id+"] - ", strings.Join(img.RepoTags, ", "), " not inserted in bucket because already exists")
		}

		bar.Increment()
	}

	bar.FinishPrint("Done.")

	return nil
}

func (r *Commands) List() {

	images, err := r.repository.GetAll()

	if err != nil {
		r.parrot.Error("Error", err)
		return
	}

	var header = models.Image{}
	var iis = [][]string{}

	for _, i := range images {
		for _, r := range i.RowsForList(r.utils) {
			iis = append(iis, r)
		}
	}

	r.parrot.Tablify(header.HeaderForList(), iis)

}

func (r *Commands) ListByNumber(co int) {
	var header = models.Image{}

	if co <= 0 {
		r.parrot.Tablify(header.HeaderForList(), nil)
		return
	}

	images, err := r.repository.GetLimit(co)

	if err != nil {
		r.parrot.Error("Error", err)
		return
	}

	var iis = [][]string{}

	for _, i := range images {
		for _, r := range i.RowsForList(r.utils) {
			iis = append(iis, r)
		}
	}

	r.parrot.Tablify(header.HeaderForList(), iis)

}

func (r *Commands) ListByName(name string) {
	imagesMap := make(map[string]bool)
	images, err := r.repository.GetAll()

	if err != nil {
		r.parrot.Error("Error", err)
		return
	}

	var header = models.Image{}
	var iis = [][]string{}

	for _, i := range images {
		for _, t := range i.Tags {
			if strings.Contains(t.ID, name) && !imagesMap[i.ID] {
				imagesMap[i.ID] = true

				for _, r := range i.RowsForList(r.utils) {
					iis = append(iis, r)
				}
			}
		}
	}

	r.parrot.Tablify(header.HeaderForList(), iis)

}

func (r *Commands) Labels() {
	images, err := r.repository.GetAll()

	if err != nil {
		r.parrot.Error("Error", err)
		return
	}

	var header = models.Image{}
	var iis = [][]string{}

	for _, i := range images {
		for _, r := range i.RowsForLabel() {
			iis = append(iis, r)
		}
	}

	r.parrot.Tablify(header.HeaderForLabel(), iis)
}

func (r *Commands) LabelsById(id string) {
	image, err := r.repository.Get(id)

	if err != nil {
		r.parrot.Error("Error", err)
		return
	}

	var header = models.Image{}
	var iis = [][]string{}

	for _, r := range image.RowsForLabel() {
		iis = append(iis, r)
	}

	r.parrot.Tablify(header.HeaderForLabel(), iis)
}

func (r *Commands) LabelsByTag(id string) {
	images, err := r.repository.FindByTag(id)

	if err != nil {
		r.parrot.Error("Error", err)
		return
	}

	var header = models.Image{}
	var iis = [][]string{}

	for _, i := range images {
		for _, r := range i.RowsForLabel() {
			iis = append(iis, r)
		}
	}

	r.parrot.Tablify(header.HeaderForLabel(), iis)
}

func (r *Commands) Info() {

	images, err := r.repository.GetAll()

	if err != nil {
		r.parrot.Error("Error", err)
		return
	}

	var header = models.Image{}
	var iis = [][]string{}

	for _, i := range images {
		for _, r := range i.RowsForInfo(r.utils) {
			iis = append(iis, r)
		}
	}

	r.parrot.Tablify(header.HeaderForInfo(), iis)
}

func (r *Commands) InfoById(id string) {
	image, err := r.repository.Get(id)

	if err != nil {
		r.parrot.Error("Error", err)
		return
	}

	var header = models.Image{}
	var iis = [][]string{}

	for _, r := range image.RowsForInfo(r.utils) {
		iis = append(iis, r)
	}

	r.parrot.Tablify(header.HeaderForInfo(), iis)
}

func (r *Commands) InfoByTag(id string) {
	images, err := r.repository.FindByTag(id)

	if err != nil {
		r.parrot.Error("Error", err)
		return
	}

	var header = models.Image{}
	var iis = [][]string{}

	for _, i := range images {
		for _, r := range i.RowsForInfo(r.utils) {
			iis = append(iis, r)
		}
	}

	r.parrot.Tablify(header.HeaderForInfo(), iis)
}

func (r *Commands) Volumes() {
	images, err := r.repository.GetAll()

	if err != nil {
		r.parrot.Error("Error", err)
		return
	}

	var header = models.Image{}
	var iis = [][]string{}

	for _, i := range images {
		for _, r := range i.RowsForVolume() {
			iis = append(iis, r)
		}
	}

	r.parrot.Tablify(header.HeaderForVolume(), iis)
}

func (r *Commands) VolumesById(id string) {
	image, err := r.repository.Get(id)

	if err != nil {
		r.parrot.Error("Error", err)
		return
	}

	var header = models.Image{}
	var iis = [][]string{}

	for _, r := range image.RowsForVolume() {
		iis = append(iis, r)
	}

	r.parrot.Tablify(header.HeaderForVolume(), iis)
}

func (r *Commands) VolumesByTag(id string) {
	images, err := r.repository.FindByTag(id)

	if err != nil {
		r.parrot.Error("Error", err)
		return
	}

	var header = models.Image{}
	var iis = [][]string{}

	for _, i := range images {
		for _, r := range i.RowsForVolume() {
			iis = append(iis, r)
		}
	}
	r.parrot.Tablify(header.HeaderForVolume(), iis)
}
