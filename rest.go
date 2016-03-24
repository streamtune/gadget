package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func serve() {
	gin.SetMode(settings.RestServerMode())

	r := gin.Default()
	r.Use(Cors())
	v1 := r.Group("api/v1")
	{
		v1.GET("/images", GetImages)
		v1.GET("/images/:repo/:tag", GetLabelByRepoTag)
		v1.GET("/images/:repo", GetLabelByTag)
		v1.GET("/labelsIndexes", GetLabelsIndexes)
		v1.POST("/labelsIndexes", GetImagesByLabel)
	}
	r.Run(":" + strconv.Itoa(settings.RestPort()))

}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func GetImages(c *gin.Context) {
	// curl -i -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/images
	var images = repository.GetAll()

	for _, img := range images {
		parrot.Debug("[" + img.ID + "] - " + strings.Join(img.RepoTags, ", ") + " [" + strconv.Itoa(len(img.Labels)) + "]")
	}

	if len(images) == 0 {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, images)
	}
}

func GetLabelByRepoTag(c *gin.Context) {
	// curl -i -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/images/
	tag := c.Param("tag")
	repo := c.Param("repo")

	var img = repository.FindByTag(repo + "/" + tag)

	parrot.Debug("[" + tag + "] - " + asJson(img.Labels) + " [" + strconv.Itoa(len(img.Labels)) + "]")

	if len(img.Labels) == 0 {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, img.Labels)
	}
}

func GetLabelByTag(c *gin.Context) {
	// curl -i -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/images/
	repo := c.Param("repo")

	var img = repository.FindByTag(repo)

	parrot.Debug("[" + repo + "] - " + asJson(img.Labels) + " [" + strconv.Itoa(len(img.Labels)) + "]")

	if len(img.Labels) == 0 {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, img.Labels)
	}
}

func GetLabelsIndexes(c *gin.Context) {
	// curl -i -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/labelsIndexes
	var labelsIndexes = repository.GetLabelsIndexes()

	for _, lbl := range labelsIndexes {
		parrot.Debug("[" + lbl.Label + "] - " + strings.Join(lbl.Ids, ", "))
	}

	if len(labelsIndexes) == 0 {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, labelsIndexes)
	}
}

func GetImagesByLabel(c *gin.Context) {
	// curl --data "lbl=vendor" -H "Content-Type: application/json" http://localhost:9080/api/v1/labelsIndexes
	lbl := c.PostForm("lbl")

	var images = repository.GetImagesByLabel(lbl)

	for _, img := range images {
		parrot.Debug("[" + img.ID + "] - " + strings.Join(img.RepoTags, ", ") + " [" + strconv.Itoa(len(img.Labels)) + "]")
	}

	if len(images) == 0 {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, images)
	}
}
