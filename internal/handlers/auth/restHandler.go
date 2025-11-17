package auth

import (
	"File_Storage/internal/services/auth"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RestHandler struct {
	service auth.Service
}

func NewRestHandler(service auth.Service) *RestHandler {
	return &RestHandler{service: service}
}

func (h *RestHandler) Register(c echo.Context) error {
	var request auth.Request

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"Invalid request": err.Error()})
	}

	if err := h.service.Register(request); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Login successfully"})
}

func (h *RestHandler) Login(c echo.Context) error {
	var request auth.Request

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	token, err := h.service.Login(request)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Register successfully", "token": token})
}
