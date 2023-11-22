package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// create a struct for database response
type Actor struct {
	ActorId   int    `json:"actor_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func findActorById(actorId string) (*Actor, error) {
	result, err := db.Query("SELECT actor_id , first_name , last_name FROM actor WHERE actor_id = ?", actorId)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	var actor Actor
	if result.Next() {
		if err := result.Scan(&actor.ActorId, &actor.FirstName, &actor.LastName); err != nil {
			return &actor, fmt.Errorf("Could not find actor id : %v", actorId)
		}
		return &actor, nil
	}

	return &actor, fmt.Errorf("could not find the resource")

}

func findActorByName(actorName string) (*[]Actor, error) {
	result, err := db.Query("SELECT actor_id , first_name , last_name FROM actor WHERE first_name like ? or last_name like ?", actorName, actorName)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	var actorList []Actor
	if result.Next() {
		var actor Actor
		if err := result.Scan(&actor.ActorId, &actor.FirstName, &actor.LastName); err != nil {
			return &actorList, fmt.Errorf("Could not find actor id : %v", actorName)
		}
		actorList = append(actorList, actor)
	}
	return &actorList, nil

}

func getActorById(c *gin.Context) {
	var actorId string = c.Param("id")
	actor, err := findActorById(actorId)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Actor not found"})
	}
	c.IndentedJSON(http.StatusOK, *actor)

}

func getActorByName(c *gin.Context) {
	var actorName string = c.Param("name")
	result, err := findActorByName(actorName)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Actor not found"})
	}
	c.IndentedJSON(http.StatusOK, *result)

}

func getActor(c *gin.Context) {
	result, err := db.Query("SELECT actor_id , first_name , last_name FROM actor limit 10")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	var actorList []Actor
	for result.Next() {
		var actor Actor
		if err := result.Scan(&actor.ActorId, &actor.FirstName, &actor.LastName); err != nil {
			fmt.Errorf("there was error while parsing actor_id : %v", actor.ActorId)
		}
		actorList = append(actorList, actor)
	}
	fmt.Println(actorList)
	c.JSON(http.StatusOK, actorList)
}

func addActor(c *gin.Context) {
	var newActor Actor

	// get the json payload from post request
	if err := c.BindJSON(&newActor); err != nil {
		log.Fatal(err.Error())
	}

	// add newActor to database
	result, err := createNewActor(&newActor)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	fmt.Println(result)
	c.JSON(http.StatusCreated, *result)

}

func createNewActor(actor *Actor) (*Actor, error) {

	// first check if the record exist in the database

	actorExist, err := db.Query("select actor_id ,first_name , last_name from  actor where first_name = ? and last_name = ?", actor.FirstName, actor.LastName)
	if err != nil {
		log.Fatal(err.Error())
	}

	if actorExist.Next() {
		if err := actorExist.Scan(&actor.ActorId, &actor.FirstName, &actor.LastName); err != nil {
			return actor, fmt.Errorf("there already exist a actor with name %s %s something went wrong with Scan", actor.FirstName, actor.LastName)
		}
		return actor, fmt.Errorf("there already exist a actor with name %s %s", actor.FirstName, actor.LastName)
	}

	_, err = db.Exec("insert into actor(first_name , last_name) values ( ? , ?)", actor.FirstName, actor.LastName)
	if err != nil {
		return nil, fmt.Errorf("problem with new add actor request, with name: %s %s", actor.FirstName, actor.LastName)
	}

	result, err := db.Query("select actor_id ,first_name , last_name from  actor where first_name = ? and last_name = ?", actor.FirstName, actor.LastName)
	if err != nil {
		return nil, fmt.Errorf("could not find the actor")
	}
	defer result.Close()

	if result.Next() {
		if err := result.Scan(&actor.ActorId, &actor.FirstName, &actor.LastName); err != nil {
			return actor, fmt.Errorf("can not locate the actor ")
		}
		return actor, nil
	}

	return actor, fmt.Errorf("could not find the resource")
}
