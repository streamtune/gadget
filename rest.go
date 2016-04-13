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
	images, err := repository.GetAll()

	if err != nil {
		c.Status(http.StatusNotAcceptable)
		return
	}

	for _, img := range images {
		parrot.Debug("[" + img.ID + "] - [" + strconv.Itoa(len(img.Tags)) + "] [" + strconv.Itoa(len(img.Labels)) + "]")
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

	images, err := repository.FindByTag(repo + "/" + tag)

	if err != nil {
		c.Status(http.StatusNotAcceptable)
		return
	}

	if len(images) == 0 {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, images)
	}

}

func GetImageByTag(c *gin.Context) {
	// curl -i -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/images/
	repo := c.Param("repo")

	images, err := repository.FindByTag(repo)

	if err != nil {
		c.Status(http.StatusNotAcceptable)
		return
	}

	if len(images) == 0 {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, images)
	}
}

func GetImagesWithLabels(c *gin.Context) {
	// curl -i -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/labels
	images, err := repository.GetImagesWithLabels()

	if err != nil {
		c.Status(http.StatusNotAcceptable)
		return
	}

	if len(images) == 0 {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, images)
	}

	c.Status(http.StatusNoContent)
}

func GetImagesWithVolumes(c *gin.Context) {
	// curl -i -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/volumes
	images, err := repository.GetImagesWithVolumes()

	if err != nil {
		c.Status(http.StatusNotAcceptable)
		return
	}

	if len(images) == 0 {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, images)
	}

	c.Status(http.StatusNoContent)
}

func GetImagesByLabel(c *gin.Context) {
	// curl --data "lbl=vendor" -H "Content-Type: application/json" http://localhost:9080/api/v1/labels
	lbl := c.PostForm("lbl")

	images, err := repository.FindByLabel(lbl)

	if err != nil {
		c.Status(http.StatusNotAcceptable)
		return
	}

	if len(images) == 0 {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, images)
	}
}

func GetImagesByVolume(c *gin.Context) {
	// curl --data "vlm=vendor" -H "Content-Type: application/json" http://localhost:9080/api/v1/labels
	vlm := c.PostForm("vlm")

	images, err := repository.FindByVolume(vlm)

	if err != nil {
		c.Status(http.StatusNotAcceptable)
		return
	}

	if len(images) == 0 {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, images)
	}

	c.Status(http.StatusNoContent)
}
