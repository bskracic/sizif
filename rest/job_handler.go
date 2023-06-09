package rest

import "C"
import (
	"github.com/bskracic/sizif/db/model"
	"github.com/bskracic/sizif/rest/dto"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Bind(main *gin.RouterGroup, h *Handler) {
	jobGroup := main.Group("/job")
	jobGroup.GET("", h.GetJobs)
	jobGroup.GET("/:id", h.GetJob)
	jobGroup.POST("", h.CreateJob)
	jobGroup.POST("/start/:id", h.StartJob)
}

func (h *Handler) GetJobs(c *gin.Context) {
	c.JSON(http.StatusOK, dto.ToJobListDto(model.RetrieveJobViews(h.Db)))
}

func (h *Handler) GetJob(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	var job model.Job
	if err := h.Db.Find(&job, id); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, job)
}

func (h *Handler) CreateJob(c *gin.Context) {

	job, err := dto.ToJob(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.Db.Preload("Runner").Create(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusCreated, job)
}

func (h *Handler) StartJob(c *gin.Context) {
	// call job start function
}
