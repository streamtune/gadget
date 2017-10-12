package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	//"github.com/emicklei/go-restful"
	"github.com/gocraft/health"
)

//should be global Var
var stream = health.NewStream()

func serve() {
	// Log to stdout!
	stream.AddSink(&health.WriterSink{os.Stdout})
	// Make sink and add it to stream
	sink := health.NewJsonPollingSink(time.Minute*5, time.Minute*20)
	stream.AddSink(sink)
	// Start the HTTP server! This will expose metrics via a JSON API.
	adr := ":5001"
	sink.StartServer(adr)

	gin.SetMode(settings.RestServerMode())

	r := gin.Default()
	r.Use(Cors())

	admin := r.Group("api/v1/admin")
	{
		admin.POST("/debug", PostDebug)
		admin.POST("/revive", PostRevive)
		admin.POST("/update", PostUpdate)
	}

	images := r.Group("api/v1/images")
	{
		images.GET("/", GetImages)
		images.GET("/:repo/:tag", GetImageByRepoTag)
		images.GET("/:repo", GetImageByTag)
		images.POST("/limit", PostImagesByLimit)
	}

	labels := r.Group("api/v1/labels")
	{
		labels.GET("/", GetImagesWithLabels)
		labels.POST("", GetImagesByLabel)
	}

	volumes := r.Group("api/v1/volumes")
	{
		volumes.GET("/", GetImagesWithVolumes)
		volumes.POST("", GetImagesByVolume)
	}
	r.Run(":" + strconv.Itoa(settings.RestPort()))

}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func PostRevive(c *gin.Context) {
	err := commands.Revive()

	if err != nil {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotAcceptable)
	}
}

func PostUpdate(c *gin.Context) {
	err := commands.Update()

	if err != nil {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotAcceptable)
	}
}

func PostDebug(c *gin.Context) {
	dbg, err := strconv.ParseBool(c.PostForm("dbg"))

	if err != nil {
		dbg = false
	}

	commands.Debug(dbg)

	c.Status(http.StatusOK)
}

func GetImages(c *gin.Context) {
	// curl -i -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/images
	images, err := repository.GetAll()

	if err != nil {
		c.Status(http.StatusNotAcceptable)
		return
	}

	if settings.DebugMode() {
		for _, img := range images {
			parrot.Debug("[" + img.ID + "] - [" + strconv.Itoa(len(img.Tags)) + "] [" + strconv.Itoa(len(img.Labels)) + "]")
		}
	}

	if len(images) == 0 {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, images)
	}
}

func PostImagesByLimit(c *gin.Context) {
	// curl -i -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/images/limit/:count
	count, err := strconv.Atoi(c.PostForm("count"))

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	images, err := repository.GetLimit(count)

	if err != nil {
		c.Status(http.StatusNotAcceptable)
		return
	}

	if settings.DebugMode() {
		for _, img := range images {
			parrot.Debug("[" + img.ID + "] - [" + strconv.Itoa(len(img.Tags)) + "] [" + strconv.Itoa(len(img.Labels)) + "]")
		}
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
