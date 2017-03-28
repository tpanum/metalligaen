package api

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/tpanum/metalligaen/models"
)

func GetTeams(c echo.Context) error {
	return c.JSON(http.StatusOK, models.Teams)
}
