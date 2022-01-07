package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Character struct {
	// Apparently `gin` requires us to set the
	// "field" names with PascalCase.
	ID    string `json:"id"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	Level int    `json:"level"`
}

func main() {
	router := gin.Default()
	router.GET("/characters", listCharacters)
	router.POST("/characters", postCharacters)
	router.GET("/characters/:id", getCharacter)

	// These endpoints aren't in the tutorial, but
	// let's create it nonetheless.
	router.PUT("/characters/:id", updateCharacter)
	router.DELETE("/characters/:id", deleteCharacter)

	router.Run("localhost:8080")
}

var characters = []Character{
	{ID: "4c0ba5d1-8139-4506-9334-08a8c3314c0d", Name: "Hades", Role: "Emet-Selch", Level: 99},
	{ID: "73d8e5e4-7a81-433f-ae87-66143e5e07b7", Name: "Venat", Role: "Former Azem", Level: 99},
	{ID: "514e0e54-9712-4207-86ba-b45d3ac1b074", Name: "Hythlodaeus", Role: "Chief of the Bureau of the Architect", Level: 99},
}

func listCharacters(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"characters": characters})
}

func postCharacters(c *gin.Context) {
	var newCharacter Character

	if err := c.BindJSON(&newCharacter); err != nil {
		return
	}

	newCharacter.ID = uuid.New().String()
	characters = append(characters, newCharacter)

	c.IndentedJSON(http.StatusCreated, newCharacter)
}

func getCharacter(c *gin.Context) {
	id := c.Param("id")

	for _, char := range characters {
		if char.ID == id {
			c.IndentedJSON(http.StatusOK, char)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Character not found"})
}

func updateCharacter(c *gin.Context) {
	var body Character

	c.ShouldBindJSON(&body)
	id := c.Param("id")

	for idx, char := range characters {
		if char.ID == id {
			// Unsure if there's a good way to do it,
			// but I just simply changed the array directly.
			// After that, I manually updated the `ID`.
			// This is so that, in case the JSON body contains
			// invalid `id`, then the valid id will still be used.
			// It's overengineering, I agree, for this case...
			characters[idx] = body
			// Ensure that ID isn't replaced.
			characters[idx].ID = id

			c.IndentedJSON(http.StatusOK, characters[idx])
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Character not found"})
}

func deleteCharacter(c *gin.Context) {
	id := c.Param("id")
	var characterIdx int = -1

	for idx, char := range characters {
		if char.ID == id {
			characterIdx = idx
			break
		}
	}

	if characterIdx == -1 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Character not found"})
		return
	}

	characters = append(characters[:characterIdx], characters[characterIdx+1:]...)
	c.IndentedJSON(http.StatusOK, gin.H{})
}
