package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Person struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	City        string `json:"city"`
	State       string `json:"state"`
	Street1     string `json:"street1"`
	Street2     string `json:"street2"`
	ZipCode     string `json:"zip_code"`
}

func HomepageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Simple Person listing API"})
}

func main() {

	db, err := sql.Open("mysql", "root:Abcd1298.@tcp(127.0.0.1:3306)/cetec")
	if err != nil {
		fmt.Print(err.Error())
	}
	defer db.Close()

	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}

	// Define router
	router := gin.Default()

	// Checking Server
	router.GET("/", HomepageHandler)

	// GET a person detail
	router.GET("/person/:id/info", func(c *gin.Context) {
		var (
			person Person
			result gin.H
		)
		id := c.Param("id")
		row := db.QueryRow("SELECT p.name, ph.number, addr.city, addr.state, addr.street1, addr.street2, addr.zip_code FROM person p INNER JOIN phone ph on ph.person_id = p.id INNER JOIN address_join addrj on addrj.person_id = p.id INNER JOIN address addr on addr.id = addrj.address_id WHERE p.id = ?;", id)
		err = row.Scan(&person.Name, &person.PhoneNumber, &person.City, &person.State, &person.Street1, &person.Street2, &person.ZipCode)
		if err != nil {
			// If no results send null
			result = gin.H{
				"result": nil,
				"count":  0,
			}
		} else {
			result = gin.H{
				"result": person,
				"count":  1,
			}
		}
		c.JSON(http.StatusOK, result)
	})

	// POST add new person
	router.POST("/person/create", func(c *gin.Context) {

		var person Person
		if err := c.ShouldBindJSON(&person); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// insert person
		stmt, err := db.Prepare("insert into person (name) values(?);")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer stmt.Close()
		res, err := stmt.Exec(person.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		personId, err := res.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// insert person info
		stmt1, err := db.Prepare("insert into phone (number, person_id) values(?, ?);")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer stmt1.Close()
		_, err = stmt1.Exec(person.PhoneNumber, personId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// insert address
		stmt2, err := db.Prepare("insert into address (city, state, street1, street2, zip_code) values(?,?,?,?,?);")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer stmt2.Close()
		res1, err := stmt2.Exec(person.City, person.State, person.Street1, person.Street2, person.ZipCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		addressId, err := res1.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// insert address join
		stmt3, err := db.Prepare("insert into address_join (person_id, address_id) values(?, ?);")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer stmt3.Close()
		_, err = stmt3.Exec(personId, addressId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("%s successfully created", person.Name),
		})

	})
	router.Run(":8083")
}
