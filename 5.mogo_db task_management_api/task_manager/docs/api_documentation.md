**API Documentation**  
🚀 This API documentation provides an overview of the various endpoints and their corresponding CRUD (Create, Read, Update, Delete) operations. It also covers how to work with variables and how to run tests against the API.
🔖 How to use this api

- Step 1: Send requestsRESTful APIs allow you to perform CRUD operations using the POST, GET, PUT, and DELETE HTTP methods.This collection contains each of these request types. Open each request and click "Send" to see what happens.
    
- Step 2: View responsesObserve the response tab for status code (200 OK), response time, and size.
    
- Step 3: Send new Body dataUpdate or add new data in "Body" in the POST request. Typically, Body data is also used in PUT request.

## Connecting to MongoDB
### To connect to MongoDB, you will need to:
    1. Install MongoDB: Download and install MongoDB on your local machine.
    2. Create a MongoDB Database: Create a new database and collection in MongoDB.
    3. Update API Connection: Update the API connection settings to point to your local MongoDB instance.

| **Route** | **Description** | **HTTP Method** | **Path** | **Request Body** | **Response** |
| --- | --- | --- | --- | --- | --- |
| **Get All Tasks** | Get a list of all tasks | GET | `/tasks` | NONE | JSON list of tasks  <br>200(OK) |
| **Get Task by ID** | Get a single task by ID | GET | `/tasks/:id` | NONE | Task object  <br>(200) |
| **Create Task** | Create a new task | POST | `/tasks` | `-- title (required)-- description (required)`  <br>`-- due_date (required)`  <br>`-- status (required) ["Pending", "in progress", "completed"]` | 201 (Created) |
| **Update Task** | Update an existing task | PUT | `/tasks/:id` | `-- title`  <br>`-- description`  <br>`-- due_date`  <br>`-- status ["Pending", "in progress", "completed"]` | 200 (OK)  <br>Task object |
| **Delete Task** | Delete a task by ID | DELETE | `/tasks/:id` | NONE | 200 (OK) |

## Get All Tasks

**Description:** Get a list of all tasks.  
**HTTP Method:** GET  
**Endpoint:** `/tasks`  
**Request Body:** None  
**Successful Response:**

- Status Code: 200 (OK)
    
- Response Body: JSON list of tasks
    

**Error Responses:**

- If there is an internal server error:
    
    - Status Code: 500 (Internal Server Error)
        
    - Response Body: `{ "error": "Internal server error" }`
        

## Get Task by ID

**Description:** Get a single task by ID.  
**HTTP Method:** GET  
**Endpoint:** `/tasks/{id}`  
**Request Body:** None  
**Successful Response:**

- Status Code: 200 (OK)
    
- Response Body: Task object
    

**Error Responses:**

- If the task with the specified ID is not found:
    
    - Status Code: 404 (Not Found)
        
    - Response Body: `{ "error": "Task not found" }`
        
- If there is an internal server error:
    
    - Status Code: 500 (Internal Server Error)
        
    - Response Body: `{ "error": "Internal server error" }`
        

## Create Task

**Description:** Create a new task.  
**HTTP Method:** POST  
**Endpoint:** `/tasks`  
**Request Body:**

jsonCopy

```
{
  "title": "Task 1",
  "description": "This is a sample task",
  "due_date": "2023-06-30",
  "status": "Pending"
}

 ```

**Successful Response:**

- Status Code: 201 (Created)
    
- Response Body: Task object
    

**Error Responses:**

- If the request body is invalid (e.g., missing required fields):
    
    - Status Code: 400 (Bad Request)
        
    - Response Body: `{ "error": "Invalid request body" }`
        
- If there is an internal server error:
    
    - Status Code: 500 (Internal Server Error)
        
    - Response Body: `{ "error": "Internal server error" }`
        

## Update Task

**Description:** Update an existing task.  
**HTTP Method:** PUT  
**Endpoint:** `/tasks/{id}`  
**Request Body:**

jsonCopy

```
{
  "title": "Updated Task 1",
  "description": "This is an updated task",
  "due_date": "2023-07-15",
  "status": "In Progress"
}

 ```

**Successful Response:**

- Status Code: 200 (OK)
    
- Response Body: Task object
    

**Error Responses:**

- If the task with the specified ID is not found:
    
    - Status Code: 404 (Not Found)
        
    - Response Body: `{ "error": "Task not found" }`
        
- If the request body is invalid (e.g., missing required fields):
    
    - Status Code: 400 (Bad Request)
        
    - Response Body: `{ "error": "Invalid request body" }`
        
- If there is an internal server error:
    
    - Status Code: 500 (Internal Server Error)
        
    - Response Body: `{ "error": "Internal server error" }`
        

## Delete Task

**Description:** Delete a task by ID.  
**HTTP Method:** DELETE  
**Endpoint:** `/tasks/{id}`  
**Request Body:** None  
**Successful Response:**

- Status Code: 200 (OK)
    

**Error Responses:**

- If the task with the specified ID is not found:
    
    - Status Code: 404 (Not Found)
        
    - Response Body: `{ "error": "Task not found" }`
        
- If there is an internal server error:
    
    - Status Code: 500 (Internal Server Error)
        
    - Response Body: `{ "error": "Internal server error" }`