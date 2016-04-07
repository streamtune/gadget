package main

import (
	"strconv"
)

// Images
type Model struct {
	ID uint64 `json:id sql:"AUTO_INCREMENT" gorm:"primary_key"`
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
	ImageID uint64 `gorm:index`
	Name    string `json:name gorm:"not null"`
	Version string `json:version`
	Tag     string `json:tag gorm:"not null"`
}

type ImageLabel struct {
	Model
	ImageID uint64 `gorm:index`
	Key     string `json:key`
	Value   string `json:value`
	Label   string `json:label`
}

func (r *Image) Header() []string {
	return []string{"ID", "CREATED", "IMAGE_NAME", "SHORT_ID", "SIZE", "VIRTUAL SIZE"}
}

func (r *Image) Rows() [][]string {
	var rs = [][]string{}

	for _, t := range r.Tags {
		rs = append(rs, []string{strconv.FormatUint(r.ID, 10), strconv.FormatInt(r.CreatedAt, 10), t.Tag, r.ShortId, r.Size, r.VirtualSize})
	}

	return rs
}
