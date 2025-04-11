package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/guardian/recipe-health-checker/models"
)

func getIndex(baseUrl string) (*models.RecipeIndex, error) {
	resp, err := http.Get(fmt.Sprintf("%s/v2/index.json", baseUrl))
	if err != nil {
		return nil, err
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		log.Printf("Server returned %d: %s", resp.StatusCode, string(content))
		return nil, errors.New("Unable to download index")
	}
	var idx models.RecipeIndex
	err = json.Unmarshal(content, &idx)
	if err != nil {
		return nil, err
	}
	return &idx, nil
}

func getRecipe(baseUrl string, i *models.IndexEntry) (*models.StructuredRecipe, error) {
	resp, err := http.Get(i.GetUrl(baseUrl))
	if err != nil {
		return nil, err
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		log.Printf("Server returned %d: %s", resp.StatusCode, string(content))
		return nil, errors.New("Unable to download content")
	}
	var recipe models.StructuredRecipe
	err = json.Unmarshal(content, &recipe)
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

func main() {
	baseUrlPtr := flag.String("base", "https://recipes.code.dev-guardianapis.com", "base URL of the recipes API to target")
	flag.Parse()

	recipeIndex, err := getIndex(*baseUrlPtr)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}

	for _, i := range recipeIndex.Recipes {
		recipe, err := getRecipe(*baseUrlPtr, i)
		if err != nil {
			log.Fatalf("%s", err)
		}
		markdown, err := recipe.RenderAsMarkdown()
		if err != nil {
			log.Fatalf("%s", err)
		}
		println(markdown)
		println("---------------------------")
		time.Sleep(time.Second * 10)
	}
}
