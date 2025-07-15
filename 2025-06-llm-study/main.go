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
	elasticIndex := flag.String("output-index", "recipe-problems", "Name of the elasticsearch index to write to")
	jsonFormat := flag.Bool("jsonFormat", false, "Set this to send the recipe as JSON format as opposed to Markdown")
	startAt := flag.Int("start", 0, "If you want to start partway through the recipes list")
	limit := flag.Int("limit", -1, "If you want to limit the number of recipes to process")
	flag.Parse()

	ai := llm.New(context.Background(), *modelName, *region)

	recipeIndex, err := getIndex(*baseUrlPtr)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}

	lastIndex := *limit
	if lastIndex == -1 || lastIndex > len(recipeIndex.Recipes) {
		lastIndex = len(recipeIndex.Recipes) - 1
	}

	if *startAt > lastIndex {
		log.Fatalf("-start value is invalid, check how many recipes there are to process!")
	}

	log.Printf("Starting at %d and processing %d recipes", *startAt, lastIndex)

	for _, i := range recipeIndex.Recipes[*startAt:lastIndex] {
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
		if err != nil {
			log.Printf("ERROR - %s", err)
			continue
		}

		result.RecipeId = i.RecipeID
		result.ComposerId = i.CapiID
		result.ModelUsed = *modelName
		spew.Dump(result)
		if !*noElastic {
			err = elasticsearch.WriteDoc(esBase, elasticIndex, result)
			if err != nil {
				log.Printf("Unable to write to Elasticsearch: %s", err)
				break
			}
		}
		println("---------------------------")
	}
}
