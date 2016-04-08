package main

import (
	"strings"
)

// Images
type Model struct {
	ID uint64 `json:id sql:"AUTO_INCREMENT" gorm:"primary_key"`
}

type Image struct {
	Model
	CreatedAt   string       `json:createdAt`
	ShortId     string       `json:shortId sql:"unique;index" gorm:"not null,unique"`
	LongId      string       `json:longId sql:"unique;index gorm:"not null,unique""`
	Labels      []ImageLabel `json:labels`
	Size        string       `json:size`
	VirtualSize string       `json:virtualSize`
	Tags        []ImageTag   `json:tags`
	Blob        ImageBlob    `json:blob`
}

type ImageBlob struct {
	ImageID uint64 `gorm:index`
	Blob    string `json:blob`
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

// Output results....

func (r *Image) HeaderForList() []string {
	return []string{"SHORT_ID", "IMAGE", "VERSION", "SIZE", "VIRTUAL_SIZE"}
}

func (r *Image) RowsForList() [][]string {
	var rs = [][]string{}

	for _, t := range r.Tags {
		//time.Unix(0, r.CreatedAt).Format("2006-01-02 15:04:05")
		rs = append(rs, []string{r.ShortId, t.Name, Truncate(t.Version), r.Size, r.VirtualSize})
	}

	return rs
}

func (r *Image) HeaderForInfo() []string {
	return []string{"SHORT_ID", "IMAGE_NAME", "ID", "SIZE", "VIRTUAL_SIZE"}
}

func (r *Image) RowsForInfo() [][]string {
	var rs = [][]string{}

	for _, t := range r.Tags {
		//time.Unix(0, r.CreatedAt).Format("2006-01-02 15:04:05")
		rs = append(rs, []string{r.ShortId, t.Tag, NameID(r.LongId), r.Size, r.VirtualSize})
	}

	return rs
}

func (r *Image) HeaderForLabel() []string {
	return []string{"SHORT_ID", "LABEL", "IMAGE_NAME"}
}

func (r *Image) RowsForLabel() [][]string {
	var rs = [][]string{}
	var tags = []string{}

	for _, t := range r.Tags {
		tags = append(tags, t.Tag)
	}

	for _, t := range r.Labels {
		//time.Unix(0, r.CreatedAt).Format("2006-01-02 15:04:05")
		rs = append(rs, []string{r.ShortId, t.Label, strings.Join(tags, ",")})
	}

	return rs
}
