package repos

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/fsouza/go-dockerclient"

	"github.com/asdine/storm"
	"github.com/pivotal-golang/bytefmt"

	"github.com/gi4nks/quant"

	models "github.com/streamtune/gadget/models"

	utils "github.com/streamtune/gadget/utils"
)

type Repository struct {
	parrot        *quant.Parrot
	configuration *utils.Configuration
	utils         *utils.Utilities

	DB *storm.DB
}

func NewRepository(p quant.Parrot, c utils.Configuration, u utils.Utilities) *Repository {
	return &Repository{parrot: &p, configuration: &c, utils: &u}
}

//
func (r *Repository) InitDB() {
	var err error

	b, err := quant.ExistsPath(r.configuration.RepositoryDirectory)
	if err != nil {
		r.parrot.Error("Got error when reading repository directory", err)
	}

	if !b {
		quant.CreatePath(r.configuration.RepositoryDirectory)
	}

	r.DB, err = storm.Open(r.configuration.RepositoryFullName(), storm.AutoIncrement())
	r.parrot.Debug("Opened database", r.configuration.RepositoryFullName)

	if err != nil {
		r.parrot.Error("Got error creating repository directory", err)
	}
}

func (r *Repository) InitSchema() {
	err := r.DB.Init(&models.Image{})

	if err != nil {
		r.parrot.Warn("Error initializing Image", err)
	}

	err = r.DB.Init(&models.ImageTag{})

	if err != nil {
		r.parrot.Warn("Error initializing ImageTag", err)
	}

	err = r.DB.Init(&models.ImageVolume{})

	if err != nil {
		r.parrot.Warn("Error initializing ImageVolume", err)
	}

	err = r.DB.Init(&models.ImageLabel{})

	if err != nil {
		r.parrot.Warn("Error initializing ImageLabel", err)
	}
}

func (r *Repository) CloseDB() {
	r.DB.Close()
}

func (r *Repository) BackupSchema() error {
	b, _ := quant.ExistsPath(r.configuration.RepositoryDirectory)
	if !b {
		return errors.New("Gadget repository path does not exist")
	}

	err := os.Rename(r.configuration.RepositoryFullName(), r.configuration.RepositoryFullName()+".bkp")

	if err != nil {
		return err
	}

	if _, err := os.Stat(r.configuration.RepositoryFullName()); err == nil {
		return os.Remove(r.configuration.RepositoryFullName())
	}

	return nil

}

// functionalities

func (r *Repository) Put(img docker.APIImages, imgDetails docker.Image) error {
	r.parrot.Debug("[" + r.utils.AsJson(img.RepoTags) + "] adding as " + r.utils.TruncateID(img.ID))

	var image = models.Image{}
	image.ID = r.utils.TruncateID(img.ID)
	image.LongId = img.ID

	image.CreatedAt = time.Unix(0, img.Created).Format("2006-01-02 15:04:05")
	image.Size = bytefmt.ByteSize(uint64(img.Size))
	image.VirtualSize = bytefmt.ByteSize(uint64(img.VirtualSize))

	err := r.DB.Save(&image)
	r.parrot.Debug("--> added image", image.ID)

	if err != nil {
		return err
	}

	image.Labels = []models.ImageLabel{}
	image.Tags = []models.ImageTag{}
	image.Volumes = []models.ImageVolume{}

	image.Blob = models.ImageBlob{}
	image.Blob.Summary = r.utils.AsJson(img)
	image.Blob.ID = image.ID

	err = r.DB.Save(&image.Blob)
	r.parrot.Debug("--> added imageBlob", image.Blob.ID)

	if err != nil {
		return err
	}

	// Adding tags
	for _, tag := range img.RepoTags {
		var imageTag = models.ImageTag{}

		err := r.DB.One("ID", tag, &imageTag)

		if err == nil {
			imageTag.ImageIDs = append(imageTag.ImageIDs, image.ID)
		} else {
			imageTag.ID = tag
			imageTag.ImageIDs = append(imageTag.ImageIDs, image.ID)
			imageTag.Name = strings.Split(tag, ":")[0]
			imageTag.Version = strings.Split(tag, ":")[1]
		}

		err = r.DB.Save(&imageTag)
		r.parrot.Debug("--> added imageTag", imageTag.ID, image.ID)

		if err != nil {
			return err
		}

		image.Tags = append(image.Tags, imageTag)
	}

	// Adding labels
	for k, v := range img.Labels {
		var imageLabel = models.ImageLabel{}

		err := r.DB.One("ID", k+":"+v, &imageLabel)

		if err == nil {
			imageLabel.ImageIDs = append(imageLabel.ImageIDs, image.ID)
		} else {
			imageLabel.ImageIDs = append(imageLabel.ImageIDs, image.ID)
			imageLabel.Key = k
			imageLabel.Value = v
			imageLabel.ID = k + ":" + v
		}

		err = r.DB.Save(&imageLabel)
		r.parrot.Debug("--> added imageLabel", imageLabel.ID, image.ID)
		if err != nil {
			return err
		}

		image.Labels = append(image.Labels, imageLabel)
	}

	// Add volumes
	image.Blob.Details = r.utils.AsJson(imgDetails)

	for k, v := range imgDetails.ContainerConfig.Volumes {
		var imageVolume = models.ImageVolume{}

		err := r.DB.One("ID", k+":"+imageVolume.Data, &imageVolume)
		if err == nil {
			imageVolume.ImageIDs = append(imageVolume.ImageIDs, image.ID)
		} else {
			imageVolume.ImageIDs = append(imageVolume.ImageIDs, image.ID)
			imageVolume.Volume = k
			imageVolume.Data = r.utils.AsJson(v)
			imageVolume.ID = k + ":" + imageVolume.Data
		}

		err = r.DB.Save(&imageVolume)
		r.parrot.Debug("--> added imageVolume", imageVolume.ID, image.ID)

		if err != nil {
			return err
		}

		image.Volumes = append(image.Volumes, imageVolume)
	}

	err = r.DB.Save(&image)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetAll() ([]models.Image, error) {
	images := []models.Image{}
	err := r.DB.All(&images)

	if err != nil {
		return nil, err
	}

	return images, err
}

func (r *Repository) GetLimit(count int) ([]models.Image, error) {
	images := []models.Image{}
	err := r.DB.All(&images)

	if err != nil {
		return nil, err
	}

	if count > len(images) {
		count = len(images)
	}

	return images[0:count], err
}

func (r *Repository) Get(id string) (models.Image, error) {
	var image = models.Image{}

	err := r.DB.One("ID", id, &image)

	if err != nil {
		return image, err
	}

	return image, err
}

func (r *Repository) Exists(id string) bool {
	var image = models.Image{}

	err := r.DB.One("ID", id, &image)

	if err != nil {
		return false
	}

	if &image == nil {
		return false
	}

	return true
}

func (r *Repository) FindByShortId(id string) (models.Image, error) {
	return r.Get(id)
}

func (r *Repository) FindByLongId(id string) (models.Image, error) {
	var image = models.Image{}

	err := r.DB.One("LongId", id, &image)

	if err != nil {
		return image, err
	}

	return image, err
}

func (r *Repository) FindByTag(tag string) ([]models.Image, error) {
	var images = []models.Image{}
	var imageTag = models.ImageTag{}

	err := r.DB.One("ID", tag, &imageTag)

	r.parrot.Debug("Tag:", tag, "ImageTag:", imageTag)

	if err != nil {
		return nil, err
	}

	if &imageTag == nil {
		r.parrot.Debug("No tag found")
		return nil, err
	}

	for _, id := range imageTag.ImageIDs {
		var image = models.Image{}

		err = r.DB.One("ID", id, &image)

		r.parrot.Debug("Found Image", r.utils.AsJson(image))
		if err != nil {
			return nil, err
		}

		images = append(images, image)
	}

	return images, nil
}

func (r *Repository) GetImagesWithLabels() ([]models.Image, error) {
	images := []models.Image{}
	imgs := []models.Image{}

	err := r.DB.All(&images)

	if err != nil {
		return nil, err
	}

	for _, img := range images {
		if len(img.Labels) > 0 {
			imgs = append(imgs, img)
		}
	}

	return imgs, err
}

func (r *Repository) FindByLabel(lbl string) ([]models.Image, error) {
	var images = []models.Image{}
	var imagesLabels = []models.ImageLabel{}

	err := r.DB.All(&imagesLabels)

	if err != nil {
		return nil, err
	}

	for _, il := range imagesLabels {
		if strings.ContainsAny(il.ID, lbl) {
			for _, id := range il.ImageIDs {
				var image = models.Image{}

				err = r.DB.One("ID", id, &image)

				r.parrot.Debug("Found Image", r.utils.AsJson(image))
				if err != nil {
					return nil, err
				}

				images = append(images, image)
			}
		}
	}

	return images, nil
}

func (r *Repository) GetImagesWithVolumes() ([]models.Image, error) {
	images := []models.Image{}
	imgs := []models.Image{}

	err := r.DB.All(&images)

	if err != nil {
		return nil, err
	}

	for _, img := range images {
		if len(img.Volumes) > 0 {
			imgs = append(imgs, img)
		}
	}

	return imgs, err
}

func (r *Repository) FindByVolume(vlm string) ([]models.Image, error) {
	var images = []models.Image{}
	var imagesVolumes = []models.ImageVolume{}

	err := r.DB.All(&imagesVolumes)

	if err != nil {
		return nil, err
	}

	for _, il := range imagesVolumes {
		if strings.ContainsAny(il.ID, vlm) {
			for _, id := range il.ImageIDs {
				var image = models.Image{}

				err = r.DB.One("ID", id, &image)

				r.parrot.Debug("Found Image", r.utils.AsJson(image))
				if err != nil {
					return nil, err
				}

				images = append(images, image)
			}
		}
	}

	return images, nil

}
