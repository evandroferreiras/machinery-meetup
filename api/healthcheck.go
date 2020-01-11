package api

import (
	"net/http"

	"github.com/labstack/echo"
)

func GetHealthCheck(c echo.Context) (err error) {
	return c.String(http.StatusOK, "OK")
}
