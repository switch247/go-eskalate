# Task Management REST API

## Introduction
This is a simple Task Management REST API built using the Go programming language and the Gin web framework. The API supports basic CRUD (Create, Read, Update, Delete) operations for managing tasks.

## Prerequisites
- Go (version 1.16 or higher)
- Gin web framework (`go get -u github.com/gin-gonic/gin`)
- MongoDB (version 3.6 or higher)

## Getting Started

1. Clone the repository:
```
git clone https://github.com/your-username/task-manager.git
```

2. Navigate to the project directory:
```
cd task-manager
```

3. Set the MongoDB connection string environment variable:
```
$env:MONGO_CONNECTION_STRING = "<replace here with your mongo connection string>"
```

4. Build the application:
```
go build -o task-manager ./main.go
```

5. Run the application:
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

type Task struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string     `json:"title"  validate:"required"`
	Description string     `json:"description"  validate:"required"`
	DueDate     time.Time  `json:"due_date"  validate:"required"`
	Status      string `json:"status"  validate:"required"`
}
```

## Validation
The Task struct uses the github.com/go-playground/validator package to validate the input data. The validate tags on the struct fields ensure that the required fields are provided and in the correct format.

Here's an example of how to use the validator:

```
// Create a new task
task := &Task{
	Description: "This is a sample task",
	DueDate:     time.Now().AddDate(0, 0, 30),
	Status:      StatusPending,
}

// Validate the task
err := validator.Validate(task)
if err != nil {
	// Handle validation errors
	fmt.Println(err)
	return
}

// The task is valid, you can proceed with creating it

```

## API Documentation
The API documentation, including detailed information about each endpoint, request and response formats, and error handling, is available in the `docs/api_documentation.md` file.

## Testing
You can use Postman or any other API testing tool to interact with the Task Management REST API. The Postman collection for this API is available in the `docs/` directory.

## Contributing
If you find any issues or have suggestions for improvements, please feel free to open an issue or submit a pull request.

## License
This project is licensed under the [MIT License](LICENSE).

Note: Make sure to replace `<replace here with your mongo connection string>` with your actual MongoDB connection string.