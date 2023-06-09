package main

import (
	"github.com/bskracic/sizif/db"
	"github.com/bskracic/sizif/db/model"
	"github.com/bskracic/sizif/rest"
	"github.com/bskracic/sizif/runtime"
	"github.com/bskracic/sizif/worker"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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

	w := worker.NewWorker(d, rnt, make(chan uint, 5))
	ticker := time.NewTicker(1 * time.Second)

	// Create a channel to receive ticks
	tick := ticker.C
	// Start a goroutine to run the function
	go func() {
		for range tick {
			w.CheckJobsToSchedule()
		}
	}()

	go func() {
		w.Start()
	}()

	panic(router.Run(":8080"))

}
