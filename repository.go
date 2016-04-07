package main

import (
	//	"os"
	"path/filepath"
	"strings"

	"github.com/fsouza/go-dockerclient"

	_ "github.com/mattn/go-sqlite3"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/pivotal-golang/bytefmt"
)

/*
const CREATE_IMAGES_TABLE_STMT string = "CREATE TABLE Images (	ID TEXT NOT NULL, Size TEXT, VirtualSize TEXT, Created TEXT, ParentID TEXT, PRIMARY KEY(ID));"
const CREATE_LABELS_TABLE_STMT string = "CREATE TABLE Labels (ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, Key TEXT, Value TEXT, ImgID TEXT NOT NULL, FOREIGN KEY(ImgID) REFERENCES Images(ID))"
const CREATE_TAGS_TABLE_STMT string = "CREATE TABLE Labels (ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, Key TEXT, Value TEXT, ImgID TEXT NOT NULL, FOREIGN KEY(ImgID) REFERENCES Images(ID))"
*/

type Repository struct {
	DB *gorm.DB
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

	r.DB, err = gorm.Open("sqlite3", repositoryFullName())
	if err != nil {
		parrot.Error("Got error creating repository directory", err)
	}
}

func (r *Repository) InitSchema() {

	if r.DB.HasTable(&ImageTag{}) {
		parrot.Debug("ImageTag already exists, removing it")
		r.DB.DropTable(&ImageTag{})
	}

	if r.DB.HasTable(&ImageLabel{}) {
		parrot.Debug("ImageLabel already exists, removing it")
		r.DB.DropTable(&ImageLabel{})
	}

	if r.DB.HasTable(&Image{}) {
		parrot.Debug("Image already exists, removing it")
		r.DB.DropTable(&Image{})
	}

	var il = &ImageLabel{}
	var it = &ImageTag{}
	var i = &Image{}

	r.DB.CreateTable(i)
	r.DB.CreateTable(il)
	r.DB.CreateTable(it)
}

func (r *Repository) CloseDB() {
	if err := r.DB.Close(); err != nil {
		parrot.Error("Error", err)
	}
}

func (r *Repository) BackupSchema() {
	b, _ := pathUtils.ExistsPath(settings.RepositoryDirectory())
	if !b {
		return
	}

	// TODO solve this
	/*
		err := os.Rename(repositoryFullName(), repositoryFullName()+".bkp")

		if err != nil {
			parrot.Error("Warning", err)
		}
	*/
}

// functionalities

func (r *Repository) Put(img docker.APIImages) {
	parrot.Debug("[" + asJson(img.RepoTags) + "] adding as " + TruncateID(img.ID))

	var image = Image{}
	image.ShortId = TruncateID(img.ID)
	image.LongId = img.ID

	image.CreatedAt = img.Created
	image.Size = bytefmt.ByteSize(uint64(img.Size))
	image.VirtualSize = bytefmt.ByteSize(uint64(img.VirtualSize))

	image.Labels = []ImageLabel{}
	image.Tags = []ImageTag{}

	image.Blob = asJson(img)

	// Adding tags
	for _, tag := range img.RepoTags {
		var imageTag = ImageTag{}

		imageTag.Name = strings.Split(tag, ":")[0]
		imageTag.Version = strings.Split(tag, ":")[1]
		imageTag.Tag = tag

		image.Tags = append(image.Tags, imageTag)
	}

	// Adding labels
	for k, v := range img.Labels {
		var imageLabel = ImageLabel{}

		imageLabel.Key = k
		imageLabel.Value = v
		imageLabel.Label = k + ":" + v

		image.Labels = append(image.Labels, imageLabel)
	}

	r.DB.Create(&image)
}

func (r *Repository) GetAll() []Image {
	images := []Image{}

	r.DB.Model(&images).Preload("Tags").Preload("Labels").Find(&images)

	return images
}

func (r *Repository) Get(id string) Image {
	var image = Image{}

	r.DB.Model(&image).Where("short_id = ?", id).Preload("Tags").Preload("Labels").First(&image)

	return image
}

func (r *Repository) Exists(id string) bool {
	var image = Image{}
	var count = -1

	r.DB.Where("short_id = ?", id).First(&image).Count(&count)

	if count == 0 {
		//parrot.Error("Error getting data", err)
		return false
	}
	parrot.Debug("Searching image with id", id, " - ", count)

	return true
}

func (r *Repository) FindByShortId(id string) Image {
	var image = Image{}

	r.DB.Model(&image).Where("short_id = ?", id).Preload("Tags").Preload("Labels").First(&image)

	return image
}

func (r *Repository) FindByLongId(id string) Image {
	var image = Image{}

	r.DB.Model(&image).Where("long_id = ?", id).Preload("Tags").Preload("Labels").First(&image)
	return image
}

func (r *Repository) FindByTag(tag string) []Image {
	var images = []Image{}
	var imageTag = ImageTag{}

	r.DB.Where("tag = ?", tag).First(&imageTag)

	if &imageTag == nil {
		parrot.Debug("No tag found")
		return nil
	}

	r.DB.Model(&images).Where("id = ?", imageTag.ImageID).Preload("Tags").Preload("Labels").Find(&images)

	return images
}

/*
func (r *Repository) GetLabelsIndexes() []LabelIndex {
	labels := []LabelIndex{}

	r.DB.View(func(tx *bolt.Tx) error {
		cc := tx.Bucket([]byte("LabelsIndex"))
		c := cc.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var label = LabelIndex{}
			err := json.Unmarshal(v, &label)
			if err != nil {
				return err
			}

			labels = append(labels, label)
		}

		return nil
	})

	return labels
}

func (r *Repository) GetImagesByLabel(lbl string) []docker.APIImages {
	images := []docker.APIImages{}

	r.DB.View(func(tx *bolt.Tx) error {
		cc := tx.Bucket([]byte("LabelsIndex"))
		c := cc.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var label = LabelIndex{}
			err := json.Unmarshal(v, &label)

			if err != nil {
				return err
			}

			if strings.Contains(label.Label, lbl) {
				for _, id := range label.Ids {
					var image = docker.APIImages{}
					cc := tx.Bucket([]byte("Images"))
					img := cc.Get([]byte(id))

					err := json.Unmarshal(img, &image)

					if err != nil {
						return err
					}
					images = append(images, image)
				}
			}
		}

		return nil
	})

	return images
}
*/

/*
func (r *Repository) FindById(id string) Command {
	var command = Command{}

	err := r.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Commands"))
		v := b.Get([]byte(id))

		err := json.Unmarshal(v, &command)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		parrot.Error("Error getting data", err)
	}

	return command
}

func (r *Repository) GetAllCommands() []Command {
	commands := []Command{}

	r.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Commands"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var command = Command{}
			err := json.Unmarshal(v, &command)
			if err != nil {
				return err
			}

			commands = append(commands, command)
		}

		return nil
	})

	return commands
}

func (r *Repository) GetLimitCommands(limit int) []Command {
	commands := []Command{}

	r.DB.View(func(tx *bolt.Tx) error {
		cc := tx.Bucket([]byte("Commands"))
		ii := tx.Bucket([]byte("CommandsIndex")).Cursor()

		var i = limit

		for k, v := ii.Last(); k != nil && i > 0; k, v = ii.Prev() {
			var command = Command{}

			parrot.Debug("--> k " + string(k) + " - v " + string(v))
			vv := cc.Get(v)

			parrot.Debug("--> vv " + string(vv))

			err := json.Unmarshal(vv, &command)
			if err != nil {
				return err
			}
			commands = append(commands, command)

			i--
		}

		return nil
	})

	return commands
}

func (r *Repository) GetExecutedCommands(count int) []ExecutedCommand {
	commands := r.GetLimitCommands(count)

	executedCommands := make([]ExecutedCommand, len(commands))

	for i := 0; i < len(commands); i++ {
		executedCommands[i] = commands[i].AsExecutedCommand(i)
	}

	return executedCommands
}
*/
