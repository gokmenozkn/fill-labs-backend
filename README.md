# Simple User Management API

This project contains a basic Go (Golang) application that serves as a simple user management API. The application is developed using the Gin web framework and utilizes the SQLite database.

## Getting Started

Clone the project to your machine and ensure that Go is installed.

1. Clone the Repository:
```bash
git clone https://github.com/gokmenozkn/fill-labs-backend.git
cd fill-labs-backend
```

2. Run the Application:
```bash
go run main.go
```

## API Endpoints
- GET /users: Get a list of all users.
- GET /users/:id: Get details of a specific user by ID.
- POST /users: Create a new user. Provide user details in the request body.
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "about": "Software Developer"
}
```
- PUT /users/:id: Update an existing user by ID. Provide updated user details in the request body.
```json
{
  "name": "Updated Name",
  "email": "updated@example.com",
  "about": "Updated Information"
}
```
- DELETE /users/:id: Delete a user by ID.

## Database
The application uses SQLite as the database, and the database file is created as `users.db` in the application's directory.

## Dependencies
- [Gin](https://github.com/gin-gonic/gin): HTTP web framework.
- [CORS](https://github.com/gin-contrib/cors): CORS middleware.
- [SQLite3 Driver](https://github.com/mattn/go-sqlite3): SQLite3 driver for Go.