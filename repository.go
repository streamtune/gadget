package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/fsouza/go-dockerclient"
)

type Repository struct {
	DB *bolt.DB
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

	r.DB, err = bolt.Open(repositoryFullName(), 0600, nil)
	if err != nil {
		parrot.Error("Got error creating repository directory", err)
	}
}

func (r *Repository) InitSchema() error {
	err := r.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Images"))
		if err != nil {
			parrot.Error("Create bucket: Images", err)
			return err
		}

		_, err := tx.CreateBucketIfNotExists([]byte("ImagesDetails"))
		if err != nil {
			parrot.Error("Create bucket: ImagesDetails", err)
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte("TagsIndex"))
		if err != nil {
			parrot.Error("Create bucket: TagsIndex", err)
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte("LabelsIndex"))
		if err != nil {
			parrot.Error("Create bucket: LabelsIndex", err)
			return err
		}

		return nil
	})

	return err
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

	err := os.Rename(repositoryFullName(), repositoryFullName()+".bkp")

	if err != nil {
		parrot.Error("Warning", err)
	}
}

// functionalities

func (r *Repository) Put(img docker.APIImages) {
	err := r.DB.Update(func(tx *bolt.Tx) error {
		cc := tx.Bucket([]byte("Images"))

		var id = TruncateID(img.ID)

		parrot.Debug("[" + asJson(img.RepoTags) + "] adding as " + id)

		encoded1, err := json.Marshal(img)

		if err != nil {
			return err
		}

		err = cc.Put([]byte(id), encoded1)

		if err != nil {
			return err
		}

		// Adding tags
		ii := tx.Bucket([]byte("TagsIndex"))

		for _, tag := range img.RepoTags {
			err = ii.Put([]byte(tag), []byte(id))

			if err != nil {
				return err
			}
		}

		// Adding Labels
		ll := tx.Bucket([]byte("LabelsIndex"))
		for k, v := range img.Labels {
			var labelIndex = LabelIndex{}

			labelIndex.Label = k + ":" + v

			parrot.Debug("[" + labelIndex.Label + "] currentLabel.")

			lbi := ll.Get([]byte(labelIndex.Label))

			if len(lbi) != 0 {
				err := json.Unmarshal(lbi, &labelIndex)
				if err != nil {
					return err
				}

				parrot.Debug("[ found ] " + asJson(labelIndex.Ids))
			}

			labelIndex.Ids = append(labelIndex.Ids, id)
			encoded1, err := json.Marshal(labelIndex)

			if err != nil {
				return err
			}

			err = ll.Put([]byte(labelIndex.Label), []byte(encoded1))

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		parrot.Error("Error inserting data", err)
	}
}

func (r *Repository) GetAll() []docker.APIImages {
	images := []docker.APIImages{}

	r.DB.View(func(tx *bolt.Tx) error {
		cc := tx.Bucket([]byte("Images"))
		c := cc.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var image = docker.APIImages{}
			err := json.Unmarshal(v, &image)
			if err != nil {
				return err
			}

			images = append(images, image)
		}

		return nil
	})

	return images
}

func (r *Repository) Get(id string) docker.APIImages {
	var image = docker.APIImages{}

	err := r.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Images"))
		v := b.Get([]byte(id))

		err := json.Unmarshal(v, &image)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		parrot.Error("Error getting data", err)
	}

	return image
}

func (r *Repository) FindByTag(tag string) docker.APIImages {
	var image = docker.APIImages{}

	err := r.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("TagsIndex"))
		v := b.Get([]byte(tag))

		if v != nil {
			cc := tx.Bucket([]byte("Images"))
			img := cc.Get([]byte(v))

			err := json.Unmarshal(img, &image)

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		parrot.Error("Error getting data", err)
	}

	return image
}

func (r *Repository) Exists(id string) bool {
	var image = docker.APIImages{}

	err := r.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Images"))
		v := b.Get([]byte(id))

		err := json.Unmarshal(v, &image)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		//parrot.Error("Error getting data", err)
		return false
	}

	return true
}

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

func extendAPIImages(slice []docker.APIImages, element docker.APIImages) []docker.APIImages {
	n := len(slice)
	if n == cap(slice) {
		// Slice is full; must grow.
		// We double its size and add 1, so if the size is zero we still grow.
		newSlice := make([]docker.APIImages, len(slice), 2*len(slice)+1)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0 : n+1]
	slice[n] = element
	return slice
}

// Append appends the items to the slice.
// First version: just loop calling Extend.
func appendAPIImages(slice []docker.APIImages, items ...docker.APIImages) []docker.APIImages {
	for _, item := range items {
		slice = extendAPIImages(slice, item)
	}
	return slice
}

func extendLabelIndex(slice []LabelIndex, element LabelIndex) []LabelIndex {
	n := len(slice)
	if n == cap(slice) {
		// Slice is full; must grow.
		// We double its size and add 1, so if the size is zero we still grow.
		newSlice := make([]LabelIndex, len(slice), 2*len(slice)+1)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0 : n+1]
	slice[n] = element
	return slice
}

// Append appends the items to the slice.
// First version: just loop calling Extend.
func appendLabelIndex(slice []LabelIndex, items ...LabelIndex) []LabelIndex {
	for _, item := range items {
		slice = extendLabelIndex(slice, item)
	}
	return slice
}
