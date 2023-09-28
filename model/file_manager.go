package model

import "time"

type Item struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	IsDirectory bool      `json:"is_directory"`
	FilePath    string    `json:"file_path"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
