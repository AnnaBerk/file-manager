package pgdb

import (
	"context"
	"database/sql"
	"errors"
	"f-manager/model"
	"f-manager/pkg/psql"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type FileManagerRepo struct {
	db *psql.Postgres
}

func NewFileManagerRepo(pg *psql.Postgres) *FileManagerRepo {
	return &FileManagerRepo{
		db: pg,
	}
}

func (r *FileManagerRepo) Rename(ctx context.Context, id string, newName string) error {
	query := `
        SELECT file_path, is_directory 
        FROM items 
        WHERE id = $1
    `
	var filePath string
	var isDirectory bool
	err := r.db.Conn.QueryRow(ctx, query, id).Scan(&filePath, &isDirectory)
	if err != nil {
		return err
	}

	// Проверка наличия файла или директории
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		return errors.New("file or directory does not exist")
	}

	// Получение нового пути
	newFilePath := filepath.Join(filepath.Dir(filePath), newName)

	if isDirectory {
		// Если это директория, обновляем пути для всех дочерних элементов
		err = r.updateChildPaths(ctx, filePath, newFilePath)
		if err != nil {
			return err
		}
	}

	// Переименование файла или папки
	err = os.Rename(filePath, newFilePath)
	if err != nil {
		return err
	}

	// Обновление пути в базе данных
	query = `
        UPDATE items 
        SET name = $2, file_path = $3
        WHERE id = $1
    `
	_, err = r.db.Conn.Exec(ctx, query, id, newName, newFilePath)
	return err
}

func (r *FileManagerRepo) CreateFolder(ctx context.Context, name string, parentID string) error {
	basePath := os.Getenv("PATH") // Установите базовый путь к вашей файловой системе
	folderPath := filepath.Join(basePath, name)
	var parentIDSql sql.NullInt64

	if parentID != "" {
		pID, err := strconv.Atoi(parentID) // преобразование parentID в int
		if err != nil {
			return err
		}
		parentFolder, err := r.GetFolderPath(ctx, parentID)
		if err != nil {
			return err
		}
		folderPath = filepath.Join(parentFolder, name)
		parentIDSql = sql.NullInt64{Int64: int64(pID), Valid: true}
	} else {
		parentIDSql = sql.NullInt64{Valid: false} // если parentID пуст, используем NULL
	}

	err := os.MkdirAll(folderPath, 0755)
	if err != nil {
		return err
	}

	query := `
        INSERT INTO items (name, is_directory, file_path, parent_id)
        VALUES ($1, true, $2, $3)
    `
	_, err = r.db.Conn.Exec(ctx, query, name, folderPath, parentIDSql)
	return err
}

func (r *FileManagerRepo) SaveFile(ctx context.Context, data []byte, filename string, parentID string) (string, error) {
	basePath := os.Getenv("PATH") // Установите базовый путь к вашей файловой системе
	filePath := filepath.Join(basePath, filename)
	var parentIDSql sql.NullInt64

	if parentID != "" {
		pID, err := strconv.Atoi(parentID) // Преобразование parentID в int
		if err != nil {
			return "", err
		}
		parentFolder, err := r.GetFolderPath(ctx, parentID)
		if err != nil {
			return "", err
		}
		filePath = filepath.Join(parentFolder, filename)
		parentIDSql = sql.NullInt64{Int64: int64(pID), Valid: true}
	} else {
		parentIDSql = sql.NullInt64{Valid: false} // Если parentID пуст, используем NULL
	}

	err := os.WriteFile(filePath, data, 0666)
	if err != nil {
		return "", err
	}

	query := `
        INSERT INTO items (name, is_directory, file_path, parent_id)
        VALUES ($1, false, $2, $3)
        RETURNING id
    `
	var id string
	err = r.db.Conn.QueryRow(ctx, query, filename, filePath, parentIDSql).Scan(&id)
	return id, err
}

func (r *FileManagerRepo) DeleteItem(ctx context.Context, id string) error {
	// Сначала получaeм информацию о файле/папке
	query := `
        SELECT file_path, is_directory 
        FROM items 
        WHERE id = $1
    `
	var filePath string
	var isDirectory bool
	err := r.db.Conn.QueryRow(ctx, query, id).Scan(&filePath, &isDirectory)
	if err != nil {
		return err
	}

	// Если это файл, удаляем его из файловой системы
	if !isDirectory {
		err = os.Remove(filePath)
		if err != nil {
			return err
		}
		// Теперь удаляем запись из базы данных
		query = `
        DELETE FROM items 
        WHERE id = $1
    `
		_, err = r.db.Conn.Exec(ctx, query, id)
		return err
	}
	return nil
}

func (r *FileManagerRepo) GetItemsByParentID(ctx context.Context, parentID string) ([]model.Item, error) {
	query := `
        SELECT id, name, is_directory, file_path, created_at, updated_at 
        FROM items 
        WHERE parent_id = $1
    `
	rows, err := r.db.Conn.Query(ctx, query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		err := rows.Scan(&item.ID, &item.Name, &item.IsDirectory, &item.FilePath, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (p *FileManagerRepo) DownloadFile(ctx context.Context, id string) ([]byte, string, error) {
	query := `
        SELECT name, file_path 
        FROM items 
        WHERE id = $1 AND is_directory = false
    `
	var name string
	var filePath string
	err := p.db.Conn.QueryRow(ctx, query, id).Scan(&name, &filePath)
	if err != nil {
		return nil, "", err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", err
	}

	return data, name, nil
}

func (r *FileManagerRepo) GetFolderPath(ctx context.Context, id string) (string, error) {
	query := `SELECT file_path FROM items WHERE id=$1`
	var path string
	err := r.db.Conn.QueryRow(ctx, query, id).Scan(&path)
	return path, err
}

func (r *FileManagerRepo) updateChildPaths(ctx context.Context, oldPath string, newPath string) error {
	query := `
		SELECT id, file_path 
		FROM items 
		WHERE file_path LIKE $1 || '%'
	`
	rows, err := r.db.Conn.Query(ctx, query, oldPath)

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var itemPath string
		err := rows.Scan(&id, &itemPath)
		if err != nil {
			return err
		}

		// Формирование нового пути
		updatedPath := strings.Replace(itemPath, oldPath, newPath, 1)

		updateQuery := `
			UPDATE items 
			SET file_path = $1
			WHERE id = $2
		`
		_, err = r.db.Conn.Exec(ctx, updateQuery, updatedPath, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *FileManagerRepo) CheckDuplicateName(ctx context.Context, name string, parentID string) (bool, error) {
	var count int
	var err error

	if parentID == "" {
		query := `
			SELECT COUNT(*)
			FROM items
			WHERE name = $1 AND parent_id IS NULL
		`
		err = r.db.Conn.QueryRow(ctx, query, name).Scan(&count)
	} else {
		query := `
			SELECT COUNT(*)
			FROM items
			WHERE name = $1 AND parent_id = $2
		`
		err = r.db.Conn.QueryRow(ctx, query, name, parentID).Scan(&count)
	}

	if err != nil {
		return false, err
	}
	return count > 0, nil
}
