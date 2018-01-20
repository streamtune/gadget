package rest

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	//"github.com/emicklei/go-restful"
	"github.com/gocraft/health"

	"github.com/gi4nks/quant"

	repos "github.com/streamtune/gadget/repos"
	utils "github.com/streamtune/gadget/utils"
)

//should be global Var
var stream = health.NewStream()

type Rest struct {
	parrot *quant.Parrot

	configuration *utils.Configuration
	repository    *repos.Repository
	commands      *repos.Commands
	utils         *utils.Utilities
}

func NewRest(p quant.Parrot, c utils.Configuration, r repos.Repository, cm repos.Commands, u utils.Utilities) *Rest {
	return &Rest{parrot: &p, configuration: &c, repository: &r, commands: &cm, utils: &u}
}

func (rs *Rest) Serve() {
	// Log to stdout!
	stream.AddSink(&health.WriterSink{os.Stdout})
	// Make sink and add it to stream
	sink := health.NewJsonPollingSink(time.Minute*5, time.Minute*20)
	stream.AddSink(sink)
	// Start the HTTP server! This will expose metrics via a JSON API.
	adr := ":5001"
	sink.StartServer(adr)

	gin.SetMode(rs.configuration.RestServerMode)

	r := gin.Default()
	r.Use(rs.cors())

	admin := r.Group("/api/v1/admin")
	{
		admin.POST("/debug", rs.postDebug)
		admin.POST("/revive", rs.postRevive)
		admin.POST("/update", rs.postUpdate)
	}

	images := r.Group("/api/v1/images")
	{
		images.GET("/", rs.getImages)
		images.GET("/:repo/:tag", rs.getImageByRepoTag)
		images.GET("/:repo", rs.getImageByTag)
		images.POST("/limit", rs.postImagesByLimit)
	}

	labels := r.Group("/api/v1/labels")
	{
		labels.GET("/", rs.getImagesWithLabels)
		labels.POST("", rs.getImagesByLabel)
	}

	volumes := r.Group("/api/v1/volumes")
	{
		volumes.GET("/", rs.getImagesWithVolumes)
		volumes.POST("", rs.getImagesByVolume)
	}

	rs.parrot.Println("Gadget Rest APIs running on port", rs.configuration.RestPort)

	r.Run(":" + strconv.Itoa(rs.configuration.RestPort))

}

func (rs *Rest) cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func (rs *Rest) postRevive(c *gin.Context) {
	err := rs.commands.Revive()

	if err != nil {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotAcceptable)
	}
}

func (rs *Rest) postUpdate(c *gin.Context) {
	// curl -i -X POST -H "Content-Type: application/json" http://localhost:9080/api/v1/admin/update

	err := rs.commands.Update()

	if err != nil {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotAcceptable)
	}
}

func (rs *Rest) postDebug(c *gin.Context) {
	dbg, err := strconv.ParseBool(c.PostForm("dbg"))

	if err != nil {
		dbg = false
	}

	rs.commands.Debug(dbg)

	c.Status(http.StatusOK)
}

func (rs *Rest) getImages(c *gin.Context) {
	// curl -i -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/images
	images, err := rs.repository.GetAll()

	rs.parrot.Println("-->", err)

	if err != nil {
		c.Status(http.StatusNotAcceptable)
		return
	}

	if rs.configuration.DebugMode {
		for _, img := range images {
			rs.parrot.Debug("[" + img.ID + "] - [" + strconv.Itoa(len(img.Tags)) + "] [" + strconv.Itoa(len(img.Labels)) + "]")
		}
	}

	if len(images) == 0 {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, images)
	}
}

func (rs *Rest) postImagesByLimit(c *gin.Context) {
	// curl -i -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/images/limit/:count
	count, err := strconv.Atoi(c.PostForm("count"))

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	images, err := rs.repository.GetLimit(count)

	if err != nil {
		c.Status(http.StatusNotAcceptable)
		return
	}

	if rs.configuration.DebugMode {
		for _, img := range images {
			rs.parrot.Debug("[" + img.ID + "] - [" + strconv.Itoa(len(img.Tags)) + "] [" + strconv.Itoa(len(img.Labels)) + "]")
		}
	}

	if len(images) == 0 {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, images)
	}
}

func (rs *Rest) getImageByRepoTag(c *gin.Context) {
	// curl -i -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/images/

	tag := c.Param("tag")
	repo := c.Param("repo")

	images, err := rs.repository.FindByTag(repo + "/" + tag)

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

func (rs *Rest) getImageByTag(c *gin.Context) {
	// curl -i -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/images/
	repo := c.Param("repo")

	images, err := rs.repository.FindByTag(repo)

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

func (rs *Rest) getImagesWithLabels(c *gin.Context) {
	// curl -i -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/labels
	images, err := rs.repository.GetImagesWithLabels()

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

func (rs *Rest) getImagesWithVolumes(c *gin.Context) {
	// curl -i -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/volumes
	images, err := rs.repository.GetImagesWithVolumes()

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

func (rs *Rest) getImagesByLabel(c *gin.Context) {
	// curl --data "lbl=vendor" -H "Content-Type: application/json" http://localhost:9080/api/v1/labels
	lbl := c.PostForm("lbl")

	images, err := rs.repository.FindByLabel(lbl)

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

func (rs *Rest) getImagesByVolume(c *gin.Context) {
	// curl --data "vlm=vendor" -H "Content-Type: application/json" http://localhost:9080/api/v1/labels
	vlm := c.PostForm("vlm")

	images, err := rs.repository.FindByVolume(vlm)

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
