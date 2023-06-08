package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogap/aop"
	"net/http"
	"reflect"
)

func (h *Handler) Before(jp aop.JoinPointer) {
	var dto any
	jp.Args().MapTo(func(c *gin.Context, t reflect.Type) {
		dtoTemp := reflect.New(t).Elem()
		if err := c.ShouldBindJSON(&dtoTemp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		dto = dtoTemp.Interface()
	})

	fmt.Printf("Before unmarshal: %s\n", dto)
}
