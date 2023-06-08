package rest

import "C"
import (
	"github.com/bskracic/sizif/rest/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Bind(main *gin.RouterGroup, h *Handler) {
	jobGroup := main.Group("/job")
	jobGroup.POST("", h.CreateJob)
}

func (h *Handler) CreateJob(c *gin.Context) {

	job, err := dto.ToJob(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.Db.Create(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusCreated, job)
}
