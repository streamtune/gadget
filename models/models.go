package models

import (
	"strconv"
	"strings"

	utils "github.com/streamtune/gadget/utils"
)

// Images
type Image struct {
	ID          string        `json:shortId storm:"id"`
	CreatedAt   string        `json:createdAt`
	LongId      string        `json:longId storm:"index"`
	Labels      []ImageLabel  `json:labels`
	Size        string        `json:size`
	VirtualSize string        `json:virtualSize`
	Tags        []ImageTag    `json:tags`
	Volumes     []ImageVolume `json:volumes`
	Blob        ImageBlob     `json:blob`
}

type ImageBlob struct {
	ID      string `json:shortId storm:"id"`
	Summary string `json:summary`
	Details string `json:details`
}

type ImageTag struct {
	ID       string   `json:tag storm:"id"`
	ImageIDs []string `json:images`
	Name     string   `json:name`
	Version  string   `json:version`
}

type ImageVolume struct {
	ID       string   `json:volumedata storm:"id"`
	ImageIDs []string `json:images`
	Volume   string   `json:volume storm:"index"`
	Data     string   `json:data`
}

type ImageLabel struct {
	ID       string   `json:label storm:"id"`
	ImageIDs []string `json:images`
	Key      string   `json:key`
	Value    string   `json:value`
}

// Output results....

func (r *Image) HeaderForList() []string {
	return []string{"SHORT_ID", "IMAGE", "VERSION", "SIZE", "VIRTUAL_SIZE", "TAGS", "LABELS", "VOLUMES"}
}

func (r *Image) RowsForList(u *utils.Utilities) [][]string {
	var rs = [][]string{}

	for _, t := range r.Tags {
		//time.Unix(0, r.CreatedAt).Format("2006-01-02 15:04:05")
		rs = append(rs, []string{r.ID, t.Name, u.Truncate(t.Version),
			r.Size, r.VirtualSize,
			strconv.Itoa(len(r.Tags)),
			strconv.Itoa(len(r.Labels)),
			strconv.Itoa(len(r.Volumes))})
	}

	return rs
}

func (r *Image) HeaderForInfo() []string {
	return []string{"SHORT_ID", "IMAGE_NAME", "ID", "SIZE", "VIRTUAL_SIZE"}
}

func (r *Image) RowsForInfo(u *utils.Utilities) [][]string {
	var rs = [][]string{}

	for _, t := range r.Tags {
		//time.Unix(0, r.CreatedAt).Format("2006-01-02 15:04:05")
		rs = append(rs, []string{r.ID, t.ID, u.NameID(r.LongId), r.Size, r.VirtualSize})
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
		tags = append(tags, t.ID)
	}

	for _, t := range r.Labels {
		//time.Unix(0, r.CreatedAt).Format("2006-01-02 15:04:05")
		rs = append(rs, []string{r.ID, t.ID, strings.Join(tags, ",")})
	}

	return rs
}

func (r *Image) HeaderForVolume() []string {
	return []string{"SHORT_ID", "VOLUME", "DATA", "IMAGE_NAME"}
}

func (r *Image) RowsForVolume() [][]string {
	var rs = [][]string{}
	var tags = []string{}

	for _, t := range r.Tags {
		tags = append(tags, t.ID)
	}

	for _, t := range r.Volumes {
		//time.Unix(0, r.CreatedAt).Format("2006-01-02 15:04:05")
		rs = append(rs, []string{r.ID, t.Volume, t.Data, strings.Join(tags, ",")})
	}

	return rs
}
