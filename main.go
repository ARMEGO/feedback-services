package main

import (
	"net/http"

	"database/sql"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

/**
1. Employee collection - id, username, rating, reviews
2. Review collection - id, rating, comments, owner, assignedBy, assignedOn, reviewedBy
*/

type feedback struct {
	ID         string `json:"id"`
	Rating     int    `json:"rating"`
	Owner      string `json:"owner"`
	Comments   string `json:"comments"`
	AssignedOn string `json:"assigned_on"`
	ReviewedBy string `json:"reviewed_by"`
}

type employee struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Rating   int    `json:"rating"`
	Reviews  int    `json:"reviews"`
}

/*
- Add employee username
*/
func insertEmployee(c *gin.Context) {
	var newEmployee employee
	if err := c.BindJSON(&newEmployee); err != nil {
		return
	}
	insertEmployee := fmt.Sprintf(`insert into employee("username") values('%s')`, newEmployee.Username)
	db := connect()
	defer db.Close()
	_, e := db.Exec(insertEmployee)
	CheckError(e)
	c.JSON(http.StatusOK, gin.H{"message": "MESSAGE_EMPLOYEE_ADDED"})
}

/*
- List employee with rating and reviews
*/
func getEmployees(c *gin.Context) {
	db := connect()
	defer db.Close()
	ratingQuery := `(SELECT COUNT(rating) FROM feedback WHERE owner = username) AS rating`
	reviewsQuery := `(SELECT COUNT(*) FROM feedback WHERE owner = username) AS reviews`
	employeeQuery := fmt.Sprintf(`SELECT employee.*, %s, %s FROM employee;`, ratingQuery, reviewsQuery)
	rows, err := db.Query(employeeQuery)
	defer rows.Close()
	CheckError(err)
	var employees []employee
	for rows.Next() {
		var empployee employee
		/**
		TODO: for overall rating, fetch the sum and divide by reviews
		*/
		err = rows.Scan(&empployee.ID, &empployee.Username, &empployee.Rating, &empployee.Reviews)
		CheckError(err)
		employees = append(employees, empployee)
	}

	CheckError(err)
	c.JSON(http.StatusOK, employees)
}

/*
- List employee feedback by username
*/
func getEmployeeFeedback(c *gin.Context) {
	id := c.Param("id")
	db := connect()
	defer db.Close()
	employeeQuery := fmt.Sprintf(`SELECT id, rating, comments, assigned_on, reviewed_by  FROM feedback WHERE owner = '%s';`, id)
	rows, err := db.Query(employeeQuery)
	defer rows.Close()
	CheckError(err)
	var employeeFeedback []feedback
	for rows.Next() {
		var empFeedback feedback
		err = rows.Scan(&empFeedback.ID, &empFeedback.Rating, &empFeedback.Comments, &empFeedback.AssignedOn, &empFeedback.ReviewedBy)
		CheckError(err)
		employeeFeedback = append(employeeFeedback, empFeedback)
	}

	CheckError(err)
	c.JSON(http.StatusOK, employeeFeedback)
}

func deleteEmployeeByID(c *gin.Context) {
	id := c.Param("id")
	deleteStmt := fmt.Sprintf(`delete from employee where id=%s`, id)
	db := connect()
	defer db.Close()
	_, e := db.Exec(deleteStmt)
	CheckError(e)
	c.JSON(http.StatusOK, gin.H{"message": "MESSAGE_EMPLOYEE_DELETED"})
}

// review CRUD starts
func insertFeedback(c *gin.Context) {
	var newFeedback feedback
	if err := c.BindJSON(&newFeedback); err != nil {
		return
	}
	insertFeedback := fmt.Sprintf(`insert into feedback("owner", "rating", "comments", "assigned_on", "reviewed_by") values('%s','%d', '%s', now(), '%s')`, newFeedback.Owner, newFeedback.Rating, newFeedback.Comments, newFeedback.ReviewedBy)
	db := connect()
	defer db.Close()
	_, e := db.Exec(insertFeedback)
	CheckError(e)

	c.JSON(http.StatusOK, gin.H{"message": "MESSAGE_FEEDBACK_ADDED"})
}

func updateFeedback(c *gin.Context) {
	id := c.Param("id")
	var newFeedback feedback
	if err := c.BindJSON(&newFeedback); err != nil {
		return
	}
	updateStmt := fmt.Sprintf(`update feedback set "rating"=%d, "comments"='%s' where "id"='%s'`, newFeedback.Rating, newFeedback.Comments, id)
	fmt.Println(updateStmt)
	db := connect()
	defer db.Close()
	_, e := db.Exec(updateStmt)
	CheckError(e)
	c.JSON(http.StatusOK, gin.H{"message": "MESSAGE_FEEDBACK_UPDATED"})
}

func main() {
	router := gin.Default()
	// - No origin allowed by default
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	router.Use(cors.New(config))
	router.POST("/employees", insertEmployee)
	router.GET("/employees", getEmployees)
	router.DELETE("/employees/:id", deleteEmployeeByID)
	router.GET("/employees/:id", getEmployeeFeedback)
	router.POST("/feedback", insertFeedback)
	router.PUT("/feedback/:id", updateFeedback)
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
