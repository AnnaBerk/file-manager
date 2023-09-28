package controller

import (
	"f-manager/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(handler *echo.Echo, service service.FileManager) {
	handler.Use(middleware.Recover())

	handler.GET("/folder/:id", GetListItems(service))
	handler.POST("/upload", UploadFile(service))
	handler.POST("/create", CreateFolder(service))
	handler.PUT("/item/:id/rename", RenameItem(service))
	handler.GET("/item/:id/download", DownloadFile(service))
	handler.DELETE("/item/:id", DeleteItem(service))
}

func GetListItems(service service.FileManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		return service.ListItems(c)
	}
}
func CreateFolder(service service.FileManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		return service.CreateFolder(c)
	}
}
func UploadFile(service service.FileManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		return service.UploadFile(c)
	}
}

func RenameItem(service service.FileManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		return service.RenameItem(c)
	}
}

func DownloadFile(service service.FileManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		return service.DownloadFile(c)
	}
}

func DeleteItem(service service.FileManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		return service.DeleteItem(c)
	}
}
