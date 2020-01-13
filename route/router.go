package route

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/evandroferreiras/machinery-meetup/api"
)

//Init the echo server
func Init() *echo.Echo{
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/healthcheck", api.GetHealthCheck)
	e.POST("/report", api.PostGenerateReport)
	e.POST("/report-consolidated", api.PostGenerateReportConsolidated)

	return e
}
