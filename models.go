package main

import (
	"strconv"
	"strings"

	"github.com/fsouza/go-dockerclient"
	"github.com/pivotal-golang/bytefmt"
)

// ImageInfo
type ImageInfo struct {
	ID     string   `json:ID`
	Tags   []string `json:tags`
	Labels int      `json:labels`
}

func (c ImageInfo) Header() []string {
	return []string{"ID", "IMAGE", "VERSION", "LABELS"}
}

func (c ImageInfo) Rows() [][]string {
	var rs = [][]string{}

	for _, t := range c.Tags {
		rs = append(rs, []string{c.ID, strings.Split(t, ":")[0], Truncate(strings.Split(t, ":")[1]), strconv.Itoa(c.Labels)})
	}

	return rs
}

func AsImageInfo(img docker.APIImages) ImageInfo {
	var i = ImageInfo{}
	i.Build(img)
	return i
}

func (ii *ImageInfo) Build(img docker.APIImages) {
	ii.ID = TruncateID(img.ID)
	ii.Tags = img.RepoTags
	ii.Labels = len(img.Labels)
}

// -----------------

// LabelIndex
type LabelIndex struct {
	Label string   `json:label`
	Ids   []string `json:ids`
}

func (c LabelIndex) Header() []string {
	return []string{"LABEL", "ID"}
}

func (c LabelIndex) Rows() [][]string {
	var rs = [][]string{}

	for _, t := range c.Ids {
		rs = append(rs, []string{c.Label, t})
	}

	return rs
}

// -----------------

// ImageLabels
type ImageLabels struct {
	ID     string            `json:id`
	Labels map[string]string `json:labels`
}

func (c ImageLabels) Header() []string {
	return []string{"ID", "LABEL"}
}

func (c ImageLabels) Rows() [][]string {
	var rs = [][]string{}

	for k, v := range c.Labels {
		rs = append(rs, []string{c.ID, "'" + k + "'" + ":'" + v + "'"})
	}

	return rs
}

func AsImageLabels(img docker.APIImages) ImageLabels {
	var i = ImageLabels{}
	i.Build(img)
	return i
}

func (ii *ImageLabels) Build(img docker.APIImages) {
	ii.ID = TruncateID(img.ID)
	ii.Labels = img.Labels
}

// -----------------

// -----------------

// ImageVolumes
type ImageVolumes struct {
	ID      string              `json:id`
	Volumes map[string]struct{} `json:volumes`
}

func (c ImageVolumes) Header() []string {
	return []string{"ID", "VOLUME"}
}

func (c ImageVolumes) Rows() [][]string {
	var rs = [][]string{}

	for k, v := range c.Volumes {
		rs = append(rs, []string{c.ID, "'" + k + "'" + ":'" + asJson(v) + "'"})
	}

	return rs
}

func AsImageVolumes(img docker.Image) ImageVolumes {
	var i = ImageVolumes{}
	i.Build(img)
	return i
}

func (ii *ImageVolumes) Build(img docker.Image) {
	ii.ID = TruncateID(img.ID)
	ii.Volumes = img.ContainerConfig.Volumes
}

// -----------------

type ImageDetail struct {
	ID          string   `json:id`
	Tags        []string `json:tags`
	Size        string   `json:size`
	VirtualSize string   `json:virtualSize`
}

func AsImageDetail(img docker.APIImages) ImageDetail {
	var i = ImageDetail{}
	i.Build(img)
	return i
}

func (ii *ImageDetail) Build(img docker.APIImages) {
	ii.ID = TruncateID(img.ID)
	ii.Size = bytefmt.ByteSize(uint64(img.Size))
	ii.VirtualSize = bytefmt.ByteSize(uint64(img.VirtualSize))
	ii.Tags = img.RepoTags
}

func (c ImageDetail) Header() []string {
	return []string{"ID", "IMAGE", "VERSION", "SIZE", "VIRTUAL SIZE"}
}

func (c ImageDetail) Rows() [][]string {
	var rs = [][]string{}

	for _, t := range c.Tags {
		rs = append(rs, []string{c.ID, strings.Split(t, ":")[0], Truncate(strings.Split(t, ":")[1]), c.Size, c.VirtualSize})
	}

	return rs
}
