package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// getAlbums function responds with the list of all albums as json

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
	//c.JSON(http.StatusOK, albums)
}

func addAlbum(c *gin.Context) {
	var newAlbum album

	// bind incoming object to album type struct data
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	//Add the new album to the slice
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, albums)
}

func getAlbumById(c *gin.Context) {
	//get id param from the request
	var id = c.Param("id")
	// find album by id
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})

}

func getEnvVariable(key string) string {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func getDataSourceName() string {
	// Load the .env file in the current directory

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var dbName = os.Getenv("DB_NAME")
	var dbPassword = os.Getenv("DB_PASSWORD")
	var dbPort = os.Getenv("DB_PORT")
	var dbAdress = os.Getenv("DB_ADDRESS")
	var dbUsername = os.Getenv("DB_USERNAME")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUsername, dbPassword, dbAdress, dbPort, dbName)
}

var db *sql.DB

// connecting a mysql database at localhost:3306/sakila  to create api for database
func init() {

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var dbName = os.Getenv("DB_NAME")
	var dbPassword = os.Getenv("DB_PASSWORD")
	var dbPort = os.Getenv("DB_PORT")
	var dbAddress = os.Getenv("DB_ADDRESS")
	var dbUsername = os.Getenv("DB_USERNAME")

	var conStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUsername, dbPassword, dbAddress, dbPort, dbName)

	// create a global connection for the database
	db, err = sql.Open("mysql", conStr)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("db connection was successful")

}

func main() {

	// set up gin router
	r := gin.Default()
	r.GET("/albums", getAlbums)
	r.POST("/albums", addAlbum)
	r.GET("/albums/:id", getAlbumById)
	// this is where tutorial end and my project begins
	r.GET("api/v1/actor", getActor)
	r.GET("api/v1/actor/:id", getActorById)
	r.GET("api/v1/actorname/:name", getActorByName)
	// Post end point to add actor
	r.POST("api/v1/actor", addActor)
	r.Run("localhost:6060")
}
