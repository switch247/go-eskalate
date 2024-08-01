Sure, here's the updated README file with the changes you requested:

# Task Management REST API

## Introduction
This is a simple Task Management REST API built using the Go programming language and the Gin web framework. The API supports basic CRUD (Create, Read, Update, Delete) operations for managing tasks.

## Prerequisites
- Go (version 1.16 or higher)
- Gin web framework (`go get -u github.com/gin-gonic/gin`)

## Getting Started

1. Clone the repository:
```
git clone https://github.com/your-username/task-manager.git
```

2. Navigate to the project directory:
```
cd task-manager
```

3. Build the application:
```
go build -o task-manager ./main.go
```

4. Run the application:
```
./task-manager
```

The API will start running on `http://localhost:8080`.

## API Endpoints

The API provides the following endpoints:

### Get all tasks
```
GET /tasks
```

### Get a specific task
```
GET /tasks/{id}
```

### Create a new task
```
POST /tasks
```
**Request body:**
```json
{
  "id": "1",
  "title": "Task 1",
  "description": "This is a sample task",
  "due_date": "2023-06-30",
  "status": "Pending"
}
```

### Update a task
```
PUT /tasks/{id}
```
**Request body:**
```json
{
  "title": "Updated Task 1",
  "description": "This is an updated task",
  "due_date": "2023-07-15",
  "status": "In Progress"
}
```

### Delete a task
```
DELETE /tasks/{id}
```

## Data Model
The `Task` struct, defined in the `models` package, has the following fields:

```go
type TaskStatus string

const (
	StatusPending   TaskStatus = "Pending"
	StatusCompleted TaskStatus = "Completed"
	StatusInProgress TaskStatus = "In Progress"
)

type Task struct {
	ID          string     `json:"id"  validate:"required"`
	Title       string     `json:"title"  validate:"required"`
	Description string     `json:"description"  validate:"required"`
	DueDate     time.Time  `json:"due_date"  validate:"required"`
	Status      TaskStatus `json:"status"  validate:"required"`
}
```

## API Documentation
The API documentation, including detailed information about each endpoint, request and response formats, and error handling, is available in the `docs/api_documentation.md` file.

## Testing
You can use Postman or any other API testing tool to interact with the Task Management REST API. The Postman collection for this API is available in the `docs/` directory.

## Contributing
If you find any issues or have suggestions for improvements, please feel free to open an issue or submit a pull request.

## License
This project is licensed under the [MIT License](LICENSE).
