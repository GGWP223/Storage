package storage

import (
	"File_Storage/internal/services/storage"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RestHandler struct {
	service storage.Service
}

func NewRestStorageServer(service storage.Service) *RestHandler {
	return &RestHandler{service: service}
}

func (h *RestHandler) GetFiles(c echo.Context) error {
	token := c.QueryParam("token")
	files, err := h.service.GetAllFiles(token)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, files)
}

func (h *RestHandler) UploadFile(c echo.Context) error {
	token := c.FormValue("token")
	fileHeader, err := c.FormFile("file")

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "File is required"})
	}

	src, err := fileHeader.Open()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":             "Failed to open file",
			"error_description": err.Error(),
		})
	}

	defer src.Close()

	data, err := io.ReadAll(src)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":             "Failed to read file",
			"error_description": err.Error(),
		})
	}

	fmt.Printf("Received upload file size: %d bytes\n", len(data))

	err = h.service.UploadFile(token, fileHeader.Filename, fileHeader.Header.Get("Content-Type"), data)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":             "Failed to upload file",
			"error_description": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "File uploaded successfully"})
}

func (h *RestHandler) DownloadFile(c echo.Context) error {
	var request storage.TokenRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Request"})
	}

	meta, data, err := h.service.GetFile(request.Token, request.FileID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":             "Failed to download file",
			"error_description": err.Error(),
		})
	}

	res := c.Response()

	res.Header().Set("Content-Disposition", `attachment; filename="`+meta.FileName+`"`)
	res.Header().Set("Content-Type", meta.MimeType)

	res.WriteHeader(http.StatusOK)

	_, err = res.Write(data)

	return err
}

func (h *RestHandler) DeleteFile(c echo.Context) error {
	var request storage.TokenRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Request"})
	}

	if err := h.service.DeleteFile(request.Token, request.FileID); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "File deleted successfully"})
}
