package api

import (
	"github.com/evandroferreiras/machinery-meetup/machinery"

	"github.com/labstack/echo"
)

type GenerateReportHandler struct {

}

func PostGenerateReport(c echo.Context) (err error) {
	type body struct {
		Language string `json:"language"`
	}
	b := new(body)
	c.Bind(&b)
	machinery.GetServer().SendGitHubTask(b.Language)
	return nil
}
