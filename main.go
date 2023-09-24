package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

/**
1. Employee collection - id, username, rating, reviews
2. Review collection - id, rating, comments, owner, assignedBy, assignedOn, reviewedBy
*/

type feedback struct {
	ID      string `json:"id"`
	Rating  int    `json:"rating"`
	Reviews int    `json:"reviews"`
}

type employee struct {
	ID       string     `json:"id"`
	Username string     `json:"username"`
	Feedback []feedback `json:"feedback"`
}

var employees = []employee{
	{ID: "e1", Username: "admin", Feedback: empFeedback},
	{ID: "e2", Username: "user", Feedback: empFeedback},
	{ID: "e3", Username: "user1", Feedback: empFeedback},
}

var empFeedback = []feedback{
	{ID: "1", Rating: 2, Reviews: 5},
	{ID: "2", Rating: 5, Reviews: 15},
	{ID: "3", Rating: 2, Reviews: 25},
	{ID: "4", Rating: 5, Reviews: 35},
	{ID: "5", Rating: 1, Reviews: 45},
	{ID: "6", Rating: 2, Reviews: 51},
	{ID: "7", Rating: 3, Reviews: 55},
	{ID: "8", Rating: 4, Reviews: 52},
	{ID: "9", Rating: 1, Reviews: 50},
}

func getEmployees(c *gin.Context) {
	db := connect()
	rows, err := db.Query(`SELECT * FROM "Employee"`)
	CheckError(err)

	defer rows.Close()
	for rows.Next() {
		var username string

		err = rows.Scan(&username)
		CheckError(err)

		fmt.Println(username)
	}

	CheckError(err)
	c.JSON(http.StatusOK, rows)
}

func getFeedbackByUser(c *gin.Context) {
	c.JSON(http.StatusOK, empFeedback)
}

func getEmployeeByID(c *gin.Context) {
	id := c.Param("id")
	for _, e := range employees {
		if e.ID == id {
			c.JSON(http.StatusOK, e)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "MESSAGE_EMPLOYEE_NOT_FOUND"})
}

func postFeedback(c *gin.Context) {
	var newFeedback feedback
	if err := c.BindJSON(&newFeedback); err != nil {
		return
	}
	empFeedback = append(empFeedback, newFeedback)
	c.JSON(http.StatusCreated, empFeedback)
}

func main() {
	router := gin.Default()
	router.GET("/employees", getEmployees)
	router.GET("/employees/:id", getEmployeeByID)
	router.GET("/feedback/:userId", getFeedbackByUser)
	router.POST("/feedback", postFeedback)
	router.Run("localhost:8080")
}

const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "root"
	dbname   = "feedbackDb"
)

func connect() *sql.DB {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// check db
	err = db.Ping()
	CheckError(err)

	fmt.Println("Connected!")
	return db
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
