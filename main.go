/*
- Create a REST API to do basic CRUD
*/
package main

import (
	"net/http"

	"database/sql"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// Employee collection
type employee struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Rating   int    `json:"rating"`
	Reviews  int    `json:"reviews"`
}

// Feedback collection
type feedback struct {
	ID         string `json:"id"`
	Rating     int    `json:"rating"`
	Owner      string `json:"owner"`
	Comments   string `json:"comments"`
	AssignedOn string `json:"assigned_on"`
	ReviewedBy string `json:"reviewed_by"`
}

/*
  - Add employee username
    @params username
*/
func insertEmployee(c *gin.Context) {
	var newEmployee employee
	// bind newEmployee to psylosd request body or throw error
	if err := c.BindJSON(&newEmployee); err != nil {
		return
	}
	insertEmployee := fmt.Sprintf(`insert into employee("username") values('%s')`, newEmployee.Username)
	db := connect()
	defer db.Close() // this makes sure db is closed afer usage
	_, e := db.Exec(insertEmployee)
	CheckError(e)
	// no error. Return response
	c.JSON(http.StatusOK, gin.H{"message": "MESSAGE_EMPLOYEE_ADDED"})
}

/*
- List employee with rating and reviews
- @params ginContext
*/
func getEmployees(c *gin.Context) {
	db := connect()
	defer db.Close()
	// TODO: get sum of ratings and divide by reviews to get overall rating. Extra feature
	ratingQuery := `(SELECT COUNT(rating) FROM feedback WHERE owner = username) AS rating`
	// Get number of reviews
	reviewsQuery := `(SELECT COUNT(*) FROM feedback WHERE owner = username) AS reviews`
	// create employee reviews query
	employeeQuery := fmt.Sprintf(`SELECT employee.*, %s, %s FROM employee;`, ratingQuery, reviewsQuery)
	rows, err := db.Query(employeeQuery)
	defer rows.Close()
	CheckError(err)
	var employees []employee // we expect results to be a list
	for rows.Next() {
		var empployee employee
		// TODO: for overall rating, fetch the sum and divide by reviews
		// assign row to employee
		err = rows.Scan(&empployee.ID, &empployee.Username, &empployee.Rating, &empployee.Reviews)
		CheckError(err)
		// finally push to list
		employees = append(employees, empployee)
	}

	CheckError(err)
	c.JSON(http.StatusOK, employees)
}

/*
- List employee feedback by username
- @params id
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

/*
- Delete employee
@params id
*/
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
	// for int use old school C programming's %d
	updateStmt := fmt.Sprintf(`update feedback set "rating"=%d, "comments"='%s' where "id"='%s'`, newFeedback.Rating, newFeedback.Comments, id)
	fmt.Println(updateStmt)
	db := connect()
	defer db.Close()
	_, e := db.Exec(updateStmt)
	CheckError(e)
	c.JSON(http.StatusOK, gin.H{"message": "MESSAGE_FEEDBACK_UPDATED"})
}

// entru point of web services
func main() {
	router := gin.Default()
	config := cors.DefaultConfig()
	// No origin allowed by default
	config.AllowOrigins = []string{"*"} // Not recommented for prod
	router.Use(cors.New(config))
	router.POST("/employees", insertEmployee)
	router.GET("/employees", getEmployees)
	router.DELETE("/employees/:id", deleteEmployeeByID)
	router.GET("/employees/:id", getEmployeeFeedback)
	router.POST("/feedback", insertFeedback)
	router.PUT("/feedback/:id", updateFeedback)
	router.Run("localhost:8080")
}

// this should be in env
const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "root"
	dbname   = "feedbackDb"
)

// TODO: Move to separate module
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

// TODO: Improve logging and error handling
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
