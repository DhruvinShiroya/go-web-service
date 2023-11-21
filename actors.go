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

func findActorById(actor_id string) (*Actor, error) {
	result, err := db.Query("SELECT actor_id , first_name , last_name FROM actor WHERE actor_id = ?", actor_id)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	var actor Actor
	if result.Next() {
		if err := result.Scan(&actor.ActorId, &actor.FirstName, &actor.LastName); err != nil {
			return &actor, fmt.Errorf("Could not find actor id : %v", actor_id)
		}
		return &actor, nil
	}

	return &actor, fmt.Errorf("could not find the resource")

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
