# POS API Documentation

## Base URL
The API is currently hosted at `http://localhost:8080`.

## Endpoints

### Health & Monitoring

#### `GET /health`
Returns the health status of the service.
- **Response**: `200 OK`
  ```json
  {
    "status": "ok"
  }
  ```

#### `GET /ping`
Simple connectivity test.
- **Response**: `200 OK`
  ```text
  server say: pong
  ```

---

### Ingredients
Resource for managing baking ingredients.

#### `GET /ingredients`
Retrieve a paginated list of ingredients.
- **Query Parameters**:
  - `page` (int, default: 1): Page number.
  - `page_size` (int, default: 20): Number of items per page.
- **Response**: `200 OK`
  ```json
  {
    "data": [
      {
        "id": 1,
        "name": "Flour",
        "unit_of_measure": "kg",
        "type": "dry",
        "stock": 100,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "page": 1,
    "page_size": 20
  }
  ```

#### `POST /ingredients`
Create a new ingredient.
- **Request Body**:
  ```json
  {
    "name": "Sugar",           // required
    "unit_of_measure": "kg",   // required
    "type": "dry",             // required
    "stock": 50                // optional, default 0
  }
  ```
- **Response**: `201 Created`
  ```json
  {
    "id": 2,
    "name": "Sugar",
    "unit_of_measure": "kg",
    "type": "dry",
    "stock": 50,
    "created_at": "2024-01-27T12:00:00Z"
  }
  ```

#### `GET /ingredients/:id`
Get a specific ingredient by ID.
- **Parameters**: `id` (int)
- **Response**: `200 OK`
  ```json
  {
    "id": 1,
    ...
  }
  ```
- **Error**: `404 Not Found`

#### `PUT /ingredients/:id`
Update an existing ingredient. Fields are optional; only provided fields are updated.
- **Parameters**: `id` (int)
- **Request Body**:
  ```json
  {
    "name": "White Sugar",
    "stock": 60
  }
  ```
- **Response**: `200 OK` (Returns the updated object)

#### `DELETE /ingredients/:id`
Delete an ingredient.
- **Parameters**: `id` (int)
- **Response**: `204 No Content`
