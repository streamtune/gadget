package main

// Images
type Model struct {
	ID uint `json:id sql:"AUTO_INCREMENT" gorm:"primary_key"`
}

type Image struct {
	Model
	CreatedAt   int64        `json:createdAt`
	ShortId     string       `json:shortId sql:"unique;index" gorm:"not null,unique"`
	LongId      string       `json:longId sql:"unique;index gorm:"not null,unique""`
	Labels      []ImageLabel `json:labels`
	Size        string       `json:size`
	VirtualSize string       `json:virtualSize`
	Tags        []ImageTag   `json:tags`
	Blob        string       `json:blob`
}

type ImageTag struct {
	Model
	ImageID uint   `gorm:index`
	Name    string `json:name gorm:"not null"`
	Version string `json:version`
}

type ImageLabel struct {
	Model
	ImageID uint   `gorm:index`
	Key     string `json:key`
	Value   string `json:value`
	Label   string `json:label`
}
