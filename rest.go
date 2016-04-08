package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func serve() {
	gin.SetMode(settings.RestServerMode())

	r := gin.Default()
	r.Use(Cors())
	v1 := r.Group("api/v1")
	{
		v1.HEAD("/revive", HeadRevive)
		v1.GET("/images", GetImages)
		v1.GET("/images/:repo/:tag", GetImageByRepoTag)
		v1.GET("/images/:repo", GetImageByTag)
		v1.GET("/labels", GetImagesWithLabels)
		v1.POST("/labels", GetImagesByLabel)
		v1.GET("/volumes", GetImagesWithVolumes)
		v1.POST("/volumes", GetImagesByVolume)
	}
	r.Run(":" + strconv.Itoa(settings.RestPort()))

}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func HeadRevive(c *gin.Context) {
	err := commands.Revive()

	if err != nil {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotAcceptable)
	}
}

func GetImages(c *gin.Context) {
	// curl -i -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/images
	var images = repository.GetAll()

	for _, img := range images {
		parrot.Debug("[" + img.ShortId + "] - [" + strconv.Itoa(len(img.Tags)) + "] [" + strconv.Itoa(len(img.Labels)) + "]")
	}

	if len(images) == 0 {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, images)
	}
}

func GetImageByRepoTag(c *gin.Context) {
	// curl -i -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/images/

	tag := c.Param("tag")
	repo := c.Param("repo")

	var img = repository.FindByTag(repo + "/" + tag)

	parrot.Debug("[" + tag + "] - " + AsJson(img.Labels) + " [" + strconv.Itoa(len(img.Labels)) + "]")

	if len(img.Labels) == 0 {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, img.Labels)
	}

}

func GetImageByTag(c *gin.Context) {
	// curl -i -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/images/
	repo := c.Param("repo")

	var img = repository.FindByTag(repo)

	parrot.Debug("[" + repo + "] - " + AsJson(img.Labels) + " [" + strconv.Itoa(len(img.Labels)) + "]")

	if len(img.Labels) == 0 {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, img.Labels)
	}
}

func GetImagesWithLabels(c *gin.Context) {
	// curl -i -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/labels
	var images = repository.GetImagesWithLabels()

	if len(images) == 0 {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, images)
	}
}

func GetImagesByLabel(c *gin.Context) {
	// curl --data "lbl=vendor" -H "Content-Type: application/json" http://localhost:9080/api/v1/labels
	lbl := c.PostForm("lbl")

	var images = repository.GetImagesByLabel(lbl)

	if len(images) == 0 {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, images)
	}
}

func GetImagesWithVolumes(c *gin.Context) {
	// curl -i -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/labels
	var images = repository.GetImagesWithVolumes()

	if len(images) == 0 {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, images)
	}
}

func GetImagesByVolume(c *gin.Context) {
	// curl --data "vlm=vendor" -H "Content-Type: application/json" http://localhost:9080/api/v1/labels
	vlm := c.PostForm("vlm")

	var images = repository.GetImagesByVolume(vlm)

	if len(images) == 0 {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, images)
	}
}
