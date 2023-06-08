package main

import (
	"github.com/bskracic/sizif/db"
	"github.com/bskracic/sizif/db/model"
	"github.com/bskracic/sizif/rest"
	"github.com/bskracic/sizif/runtime"
	"github.com/gin-gonic/gin"
	"net/http"
)

//func setAspect(h *rest.Handler) {
//	beanFactory := aop.NewClassicBeanFactory()
//	beanFactory.RegisterBean("handler", h)
//	aspect := aop.NewAspect("unmarhsal_aspect", "handler")
//	pointcut := aop.NewPointcut("pointcut_1").Execution(`CreateJob`)
//	aspect.AddPointcut(pointcut)
//	aspect.AddAdvice(&aop.Advice{Ordering: aop.Before, Method: "Before", PointcutRefID: "pointcut_1"})
//	aspect.SetBeanFactory(beanFactory)
//}

func main() {

	d := db.New()
	rnt := runtime.NewDockerRuntime()
	handler := rest.NewHandler(d, rnt)

	//setAspect(handler)

	router := gin.Default()
	apiV1 := router.Group("/api/v1/")
	rest.Bind(apiV1, handler)

	//router.Static("/admin", "./templates/static")
	router.LoadHTMLGlob("templates/*")
	router.GET("/admin", func(c *gin.Context) {
		var jobs []model.Job
		handler.Db.Preload("runner_type").Find(&jobs)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"Jobs": jobs,
		})
	})

	panic(router.Run(":8080"))
}
