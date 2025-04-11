package models

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

func TestRecipeModelDefinition(t *testing.T) {
	f, err := os.Open("fixture/mango-chow.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	var recipe StructuredRecipe
	err = json.Unmarshal(data, &recipe)
	if err != nil {
		t.Errorf("StructuredRecipe failed to unmarshal: %s", err)
		t.FailNow()
	}
}

// func TestRenderAsMarkdown(t *testing.T) {
// 	f, err := os.Open("fixture/mango-chow.json")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()
// 	data, err := io.ReadAll(f)
// 	if err != nil {
// 		panic(err)
// 	}
// 	var recipe StructuredRecipe
// 	err = json.Unmarshal(data, &recipe)
// 	if err != nil {
// 		t.Errorf("StructuredRecipe failed to unmarshal: %s", err)
// 		t.FailNow()
// 	}

// 	result, err := recipe.RenderAsMarkdown()
// 	if err == nil {
// 		println(result)
// 		t.Error("testing")
// 	} else {
// 		t.Error(err)
// 	}
// }
