package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"html/template"

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

// Initialize recipes.josn from assets.go

// func init() {
// 	recipes = make([]Recipe, 0)
// 	json.Unmarshal(Assets.Files["/recipes.json"].Data, &recipes)
// }

func loadTemplate() (*template.Template, error){
	t := template.New("")
	for name, file := range Assets.Files {
		if file.IsDir() || !strings.HasSuffix(name, ".tmpl"){
			continue
		}
		h, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		// Extract the base filename (e.g., "index.tmpl")
        baseName := name[strings.LastIndex(name, "/")+1:]
		t, err = t.New(baseName).Parse(string(h))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func StaticHandler(c *gin.Context){
	filepath := c.Param("filepath")
	data := Assets.Files["/assets"+filepath].Data
	c.Writer.Write(data)
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

	t, err := loadTemplate()
	if err != nil {
		panic(err)
	}
	router := gin.Default()
	router.SetHTMLTemplate(t)
	//router.Static("/assets", "./assets")
	//router.LoadHTMLGlob("templates/*")
	router.GET("/", IndexHandler)
	router.GET("/recipes/:id", RecipeHandler)
	router.Run()
}