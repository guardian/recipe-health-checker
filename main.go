package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/guardian/recipe-health-checker/elasticsearch"
	"github.com/guardian/recipe-health-checker/llm"
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
	modelName := flag.String("model", "", "bedrock model ID to use")
	region := flag.String("region", os.Getenv("AWS_REGION"), "AWS region to target")
	esBase := flag.String("elasticsearch", "https://localhost:8443", "Base URL for elasticsearch to stash the results")
	noElastic := flag.Bool("no-elastic", false, "Don't output to Elasticsearch")

	jsonFormat := flag.Bool("jsonFormat", false, "Set this to send the recipe as JSNO format as opposed to Markdown")
	flag.Parse()

	ai := llm.New(context.Background(), *modelName, *region)

	recipeIndex, err := getIndex(*baseUrlPtr)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}

	for _, i := range recipeIndex.RandomisedSample(5) {
		recipe, err := getRecipe(*baseUrlPtr, i)
		if err != nil {
			log.Fatalf("%s", err)
		}

		var recipeText string
		var recipeFormat llm.RecipeFormat
		if *jsonFormat {
			recipeText, err = recipe.RenderAsJson()
			recipeFormat = llm.Json
		} else {
			recipeText, err = recipe.RenderAsMarkdown()
			recipeFormat = llm.Markdown
		}
		if err != nil {
			log.Fatalf("%s", err)
		}
		fmt.Printf("%s by %s", recipe.Title, recipe.Contributors)
		//println(markdown)
		result, err := ai.RequestReview(context.Background(), recipeText, recipeFormat)
		spew.Dump(result)
		if err != nil {
			log.Printf("ERROR - %s", err)
			continue
		}
		result.RecipeId = i.RecipeID
		result.ComposerId = i.CapiID
		result.ModelUsed = *modelName
		//println(result)
		if !*noElastic {
			err = elasticsearch.WriteDoc(esBase, result)
		}
		println("---------------------------")
	}
}
