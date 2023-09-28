# FileManager Service 

FileManager Service является сервисом для управления файлами и папками. Ниже приведены детали API, предоставляемого этим сервисом.

## API Endpoints

### Получение списка элементов

- **URL**: `/folder/:id`
- **Method**: `GET`
- **URL Params**: `id=[string]`
- **Success Response**: 
  - **Code**: `200`
  - **Content**: 
    ```json
    [
      {"id": "1", "name": "File1", "type": "file"},
      {"id": "2", "name": "Folder1", "type": "folder"}
    ]
    ```

### Загрузка файла

- **URL**: `/upload`
- **Method**: `POST`
- **Data Params**: 
  - `file=[file]`
  - `parent_id=[string]`
- **Success Response**: 
  - **Code**: `200`
  - **Content**: 
    ```json
    {"id": "[string]"}
    ```

### Создание папки

- **URL**: `/create`
- **Method**: `POST`
- **Data Params**: 
  - `name=[string]`
  - `parent_id=[string]`
- **Success Response**: 
  - **Code**: `200`
  - **Content**: `"Folder created successfully"`

### Переименование элемента

- **URL**: `/item/:id/rename`
- **Method**: `PUT`
- **URL Params**: `id=[string]`
- **Data Params**: 
  - `name=[string]`
- **Success Response**: 
  - **Code**: `200`

### Загрузка файла

- **URL**: `/item/:id/download`
- **Method**: `GET`
- **URL Params**: `id=[string]`
- **Success Response**: 
  - **Code**: `200`
  - **Content**: `file`

### Удаление элемента

- **URL**: `/item/:id`
- **Method**: `DELETE`
- **URL Params**: `id=[string]`
- **Success Response**: 
  - **Code**: `200`


