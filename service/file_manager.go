package service

import (
	"context"
	"f-manager/repo"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

type FileManager interface {
	ListItems(c echo.Context) error
	CreateFolder(c echo.Context) error
	UploadFile(c echo.Context) error
	RenameItem(c echo.Context) error
	DownloadFile(c echo.Context) error
	DeleteItem(c echo.Context) error
}

type FileManagerService struct {
	FileManagerRepo repo.FileManagerRepo
}

func NewFileManagerService(FileManagerRepo repo.FileManagerRepo) FileManager {
	return &FileManagerService{FileManagerRepo: FileManagerRepo}
}

func (fms *FileManagerService) ListItems(c echo.Context) error {
	parentID := c.Param("id")

	items, err := fms.FileManagerRepo.GetItemsByParentID(context.Background(), parentID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, items)
}

func (fms *FileManagerService) CreateFolder(c echo.Context) error {
	var folder struct {
		Name     string `json:"name"`
		ParentID string `json:"parent_id"`
	}

	if err := c.Bind(&folder); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Проверка на дубликат имени
	exists, err := fms.FileManagerRepo.CheckDuplicateName(context.Background(), folder.Name, folder.ParentID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if exists {
		return c.JSON(http.StatusBadRequest, "A folder with this name already exists")
	}

	err = fms.FileManagerRepo.CreateFolder(context.Background(), folder.Name, folder.ParentID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Folder created successfully")
}

func (fms *FileManagerService) UploadFile(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to read file")
	}

	parentID := c.FormValue("parent_id")

	// Проверка на дубликат имени
	exists, err := fms.FileManagerRepo.CheckDuplicateName(context.Background(), file.Filename, parentID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if exists {
		return echo.NewHTTPError(http.StatusBadRequest, "A file with this name already exists")
	}

	src, err := file.Open()
	defer src.Close()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to open file")
	}

	data, err := io.ReadAll(src)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read file data")
	}

	id, err := fms.FileManagerRepo.SaveFile(context.Background(), data, file.Filename, parentID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save file")
	}

	return c.JSON(http.StatusOK, map[string]string{"id": id})
}

func (fms *FileManagerService) RenameItem(c echo.Context) error {
	id := c.Param("id")
	var request struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	// Проверка на дубликаты имен (для ускорения и уменьщения обращений к бд эту информацию можно было бы хранить в кэше)
	exists, err := fms.FileManagerRepo.CheckDuplicateName(context.Background(), id, request.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check duplicate names")
	}
	if exists {
		return echo.NewHTTPError(http.StatusBadRequest, "Duplicate name")
	}

	err = fms.FileManagerRepo.Rename(context.Background(), id, request.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to rename item")
	}

	return c.NoContent(http.StatusOK)
}

func (fms *FileManagerService) DownloadFile(c echo.Context) error {
	id := c.Param("id")

	data, filename, err := fms.FileManagerRepo.DownloadFile(context.Background(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve file")
	}

	// Установка заголовка Content-Disposition для указания имени файла
	c.Response().Header().Set("Content-Disposition", "attachment; filename="+filename)
	// Отправка файла
	return c.Blob(http.StatusOK, "application/octet-stream", data)
}

func (fms *FileManagerService) DeleteItem(c echo.Context) error {
	id := c.Param("id")

	err := fms.FileManagerRepo.DeleteItem(context.Background(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete item")
	}

	return c.NoContent(http.StatusOK)
}
