package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsouza/go-dockerclient"

	"github.com/asdine/storm"
	"github.com/pivotal-golang/bytefmt"
)

type Repository struct {
	DB *storm.DB
}

// HELPER FUNCTIONS
func repositoryFullName() string {
	return settings.RepositoryDirectory() + string(filepath.Separator) + settings.RepositoryFile()
}

//
func (r *Repository) InitDB() {
	var err error

	b, err := pathUtils.ExistsPath(settings.RepositoryDirectory())
	if err != nil {
		parrot.Error("Got error when reading repository directory", err)
	}

	if !b {
		pathUtils.CreatePath(settings.RepositoryDirectory())
	}

	r.DB, err = storm.Open(repositoryFullName(), storm.AutoIncrement())
	parrot.Debug("Opened database", repositoryFullName())

	if err != nil {
		parrot.Error("Got error creating repository directory", err)
	}
}

func (r *Repository) InitSchema() {
	err := r.DB.Init(&Image{})

	if err != nil {
		parrot.Warn("Error initializing Image", err)
	}

	err = r.DB.Init(&ImageTag{})

	if err != nil {
		parrot.Warn("Error initializing ImageTag", err)
	}

	err = r.DB.Init(&ImageVolume{})

	if err != nil {
		parrot.Warn("Error initializing ImageVolume", err)
	}

	err = r.DB.Init(&ImageLabel{})

	if err != nil {
		parrot.Warn("Error initializing ImageLabel", err)
	}
}

func (r *Repository) CloseDB() {
	r.DB.Close()
}

func (r *Repository) BackupSchema() error {
	b, _ := pathUtils.ExistsPath(settings.RepositoryDirectory())
	if !b {
		return errors.New("Gadget repository path does not exist")
	}

	err := os.Rename(repositoryFullName(), repositoryFullName()+".bkp")

	if err != nil {
		return err
	}

	if _, err := os.Stat(repositoryFullName()); err == nil {
		return os.Remove(repositoryFullName())
	}

	return nil

}

// functionalities

func (r *Repository) Put(img docker.APIImages, imgDetails docker.Image) error {
	parrot.Debug("[" + AsJson(img.RepoTags) + "] adding as " + TruncateID(img.ID))

	var image = Image{}
	image.ID = TruncateID(img.ID)
	image.LongId = img.ID

	image.CreatedAt = time.Unix(0, img.Created).Format("2006-01-02 15:04:05")
	image.Size = bytefmt.ByteSize(uint64(img.Size))
	image.VirtualSize = bytefmt.ByteSize(uint64(img.VirtualSize))

	err := r.DB.Save(&image)
	parrot.Debug("--> added image", image.ID)

	if err != nil {
		return err
	}

	image.Labels = []ImageLabel{}
	image.Tags = []ImageTag{}
	image.Volumes = []ImageVolume{}

	image.Blob = ImageBlob{}
	image.Blob.Summary = AsJson(img)
	image.Blob.ID = image.ID

	err = r.DB.Save(&image.Blob)
	parrot.Debug("--> added imageBlob", image.Blob.ID)

	if err != nil {
		return err
	}

	// Adding tags
	for _, tag := range img.RepoTags {
		var imageTag = ImageTag{}

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
		parrot.Debug("--> added imageTag", imageTag.ID, image.ID)

		if err != nil {
			return err
		}

		image.Tags = append(image.Tags, imageTag)
	}

	// Adding labels
	for k, v := range img.Labels {
		var imageLabel = ImageLabel{}

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
		parrot.Debug("--> added imageLabel", imageLabel.ID, image.ID)
		if err != nil {
			return err
		}

		image.Labels = append(image.Labels, imageLabel)
	}

	// Add volumes
	image.Blob.Details = AsJson(imgDetails)

	for k, v := range imgDetails.ContainerConfig.Volumes {
		var imageVolume = ImageVolume{}

		err := r.DB.One("ID", k+":"+imageVolume.Data, &imageVolume)
		if err == nil {
			imageVolume.ImageIDs = append(imageVolume.ImageIDs, image.ID)
		} else {
			imageVolume.ImageIDs = append(imageVolume.ImageIDs, image.ID)
			imageVolume.Volume = k
			imageVolume.Data = AsJson(v)
			imageVolume.ID = k + ":" + imageVolume.Data
		}

		err = r.DB.Save(&imageVolume)
		parrot.Debug("--> added imageVolume", imageVolume.ID, image.ID)

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

func (r *Repository) GetAll() ([]Image, error) {
	images := []Image{}
	err := r.DB.All(&images)

	if err != nil {
		return nil, err
	}

	return images, err
}

func (r *Repository) GetLimit(count int) ([]Image, error) {
	images := []Image{}
	err := r.DB.All(&images)

	if err != nil {
		return nil, err
	}

	if count > len(images) {
		count = len(images)
	}

	return images[0:count], err
}

func (r *Repository) Get(id string) (Image, error) {
	var image = Image{}

	err := r.DB.One("ID", id, &image)

	if err != nil {
		return image, err
	}

	return image, err
}

func (r *Repository) Exists(id string) bool {
	var image = Image{}

	err := r.DB.One("ID", id, &image)

	if err != nil {
		return false
	}

	if &image == nil {
		return false
	}

	return true
}

func (r *Repository) FindByShortId(id string) (Image, error) {
	return r.Get(id)
}

func (r *Repository) FindByLongId(id string) (Image, error) {
	var image = Image{}

	err := r.DB.One("LongId", id, &image)

	if err != nil {
		return image, err
	}

	return image, err
}

func (r *Repository) FindByTag(tag string) ([]Image, error) {
	var images = []Image{}
	var imageTag = ImageTag{}

	err := r.DB.One("ID", tag, &imageTag)

	parrot.Debug("Tag:", tag, "ImageTag:", imageTag)

	if err != nil {
		return nil, err
	}

	if &imageTag == nil {
		parrot.Debug("No tag found")
		return nil, err
	}

	for _, id := range imageTag.ImageIDs {
		var image = Image{}

		err = r.DB.One("ID", id, &image)

		parrot.Debug("Found Image", AsJson(image))
		if err != nil {
			return nil, err
		}

		images = append(images, image)
	}

	return images, nil
}

func (r *Repository) GetImagesWithLabels() ([]Image, error) {
	images := []Image{}
	imgs := []Image{}

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

func (r *Repository) FindByLabel(lbl string) ([]Image, error) {
	var images = []Image{}
	var imagesLabels = []ImageLabel{}

	err := r.DB.All(&imagesLabels)

	if err != nil {
		return nil, err
	}

	for _, il := range imagesLabels {
		if strings.ContainsAny(il.ID, lbl) {
			for _, id := range il.ImageIDs {
				var image = Image{}

				err = r.DB.One("ID", id, &image)

				parrot.Debug("Found Image", AsJson(image))
				if err != nil {
					return nil, err
				}

				images = append(images, image)
			}
		}
	}

	return images, nil
}

func (r *Repository) GetImagesWithVolumes() ([]Image, error) {
	images := []Image{}
	imgs := []Image{}

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

func (r *Repository) FindByVolume(vlm string) ([]Image, error) {
	var images = []Image{}
	var imagesVolumes = []ImageVolume{}

	err := r.DB.All(&imagesVolumes)

	if err != nil {
		return nil, err
	}

	for _, il := range imagesVolumes {
		if strings.ContainsAny(il.ID, vlm) {
			for _, id := range il.ImageIDs {
				var image = Image{}

				err = r.DB.One("ID", id, &image)

				parrot.Debug("Found Image", AsJson(image))
				if err != nil {
					return nil, err
				}

				images = append(images, image)
			}
		}
	}

	return images, nil

}
