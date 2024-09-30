# FileManager Service 

FileManager Service allows to manage files and folders.

## API Endpoints

### Get elements list

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

### File upload

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

### Create folder

- **URL**: `/create`
- **Method**: `POST`
- **Data Params**: 
  - `name=[string]`
  - `parent_id=[string]`
- **Success Response**: 
  - **Code**: `200`
  - **Content**: `"Folder created successfully"`

### Rename element

- **URL**: `/item/:id/rename`
- **Method**: `PUT`
- **URL Params**: `id=[string]`
- **Data Params**: 
  - `name=[string]`
- **Success Response**: 
  - **Code**: `200`

### File download

- **URL**: `/item/:id/download`
- **Method**: `GET`
- **URL Params**: `id=[string]`
- **Success Response**: 
  - **Code**: `200`
  - **Content**: `file`

### Delete element

- **URL**: `/item/:id`
- **Method**: `DELETE`
- **URL Params**: `id=[string]`
- **Success Response**: 
  - **Code**: `200`


