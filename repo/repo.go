package repo

import (
	"context"
	"f-manager/model"
)

type FileManagerRepo interface {
	Rename(ctx context.Context, id string, newName string) error
	CreateFolder(ctx context.Context, name string, parentID string) error
	SaveFile(ctx context.Context, data []byte, filename string, parentID string) (string, error)
	DeleteItem(ctx context.Context, id string) error
	GetItemsByParentID(ctx context.Context, parentID string) ([]model.Item, error)
	DownloadFile(ctx context.Context, id string) ([]byte, string, error)
	GetFolderPath(ctx context.Context, id string) (string, error)
	CheckDuplicateName(ctx context.Context, name string, parentID string) (bool, error)
}
