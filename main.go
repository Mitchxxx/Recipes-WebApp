package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Recipe struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Ingredients []Ingredients `json:"ingredients"`
	Steps []string `json:"steps"`
	ImageUrl string `json:"imageURL"`
}

type Ingredients struct {
	Quantity string `json:"quantity"`
	Name string `json:"name"`
	Type string `json:"type"`
}
var recipes []Recipe

func init (){
	recipes = make([]Recipe, 0)
	file, err := os.Open("recipes.json")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	   	return
	}
	defer file.Close()
	fileData, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	   	return
	}
	if err = json.Unmarshal(fileData, &recipes); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
		return
	}

}

func IndexHandler(c *gin.Context) {
	//c.File("index.html")
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"recipes": recipes,
	})

}

func RecipeHandler(c *gin.Context) {
	for _, recipe := range recipes {
		if recipe.ID == c.Param("id") {
			c.HTML(http.StatusOK, "recipe.tmpl", gin.H{
				"recipe": recipe,
			})
			return
		}
	}
	c.File("404.Html")
}

func main () {
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*")
	router.GET("/", IndexHandler)
	router.GET("/recipes/:id", RecipeHandler)
	router.Run()
}