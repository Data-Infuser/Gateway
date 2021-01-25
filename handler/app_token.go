package handler

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/promptech1/infuser-gateway/models"
)

func (h *Handler) RegistAppToken(c echo.Context) error {
	c.Request().Context()
	appToken := &models.AppToken{}
	if err := c.Bind(appToken); err != nil {
		return err
	}

	return nil
}
