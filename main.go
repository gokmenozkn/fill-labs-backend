package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	About     string    `json:"about"`
}

func main() {
	// Connect to the SQLite database
	var err error
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the user table in the database
	createTable()

	// Create a Gin router
	router := gin.Default()

	// Enable CORS middleware
	router.Use(cors.Default())

	// Define REST endpoints
	router.GET("/users", getUsers)
	router.GET("/users/:id", getUser)
	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)

	// Run the application on port 8080
	router.Run(":8080")
}

func createTable() {
	// Create the 'users' table if it doesn't exist
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			email TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			about TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func getUsers(c *gin.Context) {
	// Retrieve all users from the database
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.About)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

func getUser(c *gin.Context) {
	// Retrieve a specific user by ID from the database
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user User
	err = db.QueryRow("SELECT * FROM users WHERE id=?", id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.About)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func createUser(c *gin.Context) {
	// Create a new user based on JSON input
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert the new user into the database
	result, err := db.Exec("INSERT INTO users (name, email, about) VALUES (?, ?, ?)",
		newUser.Name, newUser.Email, newUser.About)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the last inserted ID and update the newUser struct
	userID, _ := result.LastInsertId()
	newUser.ID = int(userID)

	// Return the newly created user as JSON
	c.JSON(http.StatusCreated, newUser)
}

func updateUser(c *gin.Context) {
	// Update an existing user by ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Retrieve updated user data from JSON input
	var updatedUser User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Execute the SQL query to update the user in the database
	_, err = db.Exec("UPDATE users SET name=?, email=?, about=? WHERE id=?", updatedUser.Name, updatedUser.Email, updatedUser.About, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set the updated user ID and return it as JSON
	updatedUser.ID = id
	c.JSON(http.StatusOK, updatedUser)
}

func deleteUser(c *gin.Context) {
	// Delete a user by ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Execute the SQL query to delete the user from the database
	_, err = db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a success message as JSON
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
