package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	conf "github.com/praveennagaraj97/shoppers-gocommerce/config"
	"github.com/praveennagaraj97/shoppers-gocommerce/middlewares"
	"github.com/praveennagaraj97/shoppers-gocommerce/pkg/color"
	logger "github.com/praveennagaraj97/shoppers-gocommerce/pkg/log"
	"github.com/praveennagaraj97/shoppers-gocommerce/repository"
)

type Router struct {
	conf        *conf.GlobalConfiguration // Global Configuration
	engine      *gin.Engine
	middlewares *middlewares.Middlewares
	repos       *repository.Repositories
}

func (r *Router) InitializeRouter(cfg *conf.GlobalConfiguration) {
	r.conf = cfg
	r.listenAndSeve()
}

// ListenAndServe starts the REST API.
func (router *Router) listenAndSeve() {
	if router.conf.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// initialize middlewares
	middlewares := middlewares.Middlewares{}
	middlewares.Initialize(router.conf)
	router.middlewares = &middlewares

	// initialize repo for router controllers
	router.initializeRepositories()

	r := gin.New()
	// provide gin engine to all incoming routes
	router.engine = r
	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true
	r.SetTrustedProxies([]string{router.conf.Domain})
	r.Use(gin.Logger())
	r.Use(router.corsMiddleware())

	r.Static("/static", "./static")
	r.Use(middlewares.AcceptLanguageParser())

	// list Routes
	router.globalRoutes()
	router.userRoutes()
	router.categoriesRoutes()
	router.assetsRoutes()

	// 404
	r.Use(func(ctx *gin.Context) {
		ctx.JSON(404, map[string]interface{}{
			"message": "Route not available",
		})
	})

	if err := r.Run(fmt.Sprintf(":%s", router.conf.Port)); err != nil {
		logger.PrintLog("Failed to start server", color.Red)
		logger.ErrorLogFatal(err)
	}
}

var trustedDomains map[string]bool = map[string]bool{
	"http://localhost:8080": true,
	"http://localhost:3200": true,
}

func (r *Router) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")

		origin := c.Request.Header.Get("Origin")
		c.Writer.Header().Set("Server", "api-shopee-go")

		if trustedDomains[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func (r *Router) initializeRepositories() {
	repos := &repository.Repositories{}
	repos.Initialize(r.conf.DatabaseClient, r.conf.DatabaseName)

	r.repos = repos
}
