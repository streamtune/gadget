package main

// Images
type Model struct {
	ID uint `json:id`
}

type Image struct {
	Model
	CreatedAt   int64        `json:createdAt`
	ImageId     string       `json:imageId`
	ImageLongId string       `json:imageLongId`
	Labels      []ImageLabel `json:labels`
	Size        string       `json:size`
	VirtualSize string       `json:virtualSize`
	Tags        []ImageTag   `json:tags`
	Blob        string       `json:blob`
}

type ImageTag struct {
	Model
	ImageId uint   `gorm:"index"`
	Name    string `json:name`
	Version string `json:version`
}

type ImageLabel struct {
	Model
	ImageId uint   `gorm:"index"`
	Key     string `json:key`
	Value   string `json:value`
	Label   string `json:label`
}
