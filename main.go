package main

import (
	"github.com/bskracic/sizif/db"
	"github.com/bskracic/sizif/db/model"
	"github.com/bskracic/sizif/rest"
	"github.com/bskracic/sizif/runtime"
	"github.com/bskracic/sizif/worker"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"
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

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sizif_http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status", "duration"},
	)
	jobsRunTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "sizif_jobs_run_total",
			Help: "Total number of jobs run",
		})
	jobsFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "sizif_jobs_failed_total",
			Help: "Total number of jobs failed",
		})
)

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

	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(jobsRunTotal)
	prometheus.MustRegister(jobsFailedTotal)
	router.Use(PrometheusMiddleware())
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	w := worker.NewWorker(d, rnt, make(chan uint, 5))
	w.TotalJobsRun = jobsRunTotal
	w.TotalJobsFailed = jobsFailedTotal

	ticker := time.NewTicker(1 * time.Second)
	tick := ticker.C
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

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start a timer to measure request duration
		start := time.Now()

		// Process the request
		c.Next()

		// Calculate the request duration
		duration := time.Since(start)

		// Increment the HTTP requests counter
		httpRequestsTotal.With(prometheus.Labels{
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"status":   strconv.Itoa(c.Writer.Status()),
			"duration": duration.String(),
		}).Inc()
	}
}
