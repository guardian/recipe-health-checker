package models

import (
	"fmt"
	"strings"
	"text/template"
)

type RecipeImage struct {
	URL          string  `json:"url"`
	MediaID      string  `json:"mediaId"`
	CropID       string  `json:"cropId"`
	Source       *string `json:"source,omitempty"`
	Photographer *string `json:"photographer,omitempty"`
	ImageType    *string `json:"imageType,omitempty"`
	Caption      *string `json:"caption,omitempty"`
	MediaAPIURI  *string `json:"mediaApiUri,omitempty"`
	Width        int     `json:"width"`
	Height       int     `json:"height"`
}

type AmountRange struct {
	Min float32 `json:"min"`
	Max float32 `json:"max"`
}

func formatNumber(n float32) string {
	if n == float32(int64(n)) {
		return fmt.Sprintf("%.0f", n)
	} else {
		return fmt.Sprintf("%.2f", n)
	}
}

func (a AmountRange) Render() string {

	if a.Min == a.Max {
		if a.Min == 0 {
			return ""
		} else {
			return formatNumber(a.Max)
		}
	} else {
		return fmt.Sprintf("%s - %s", formatNumber(a.Min), formatNumber(a.Max))
	}
}

type IngredientData struct {
	Name         string      `json:"name"`
	IngredientID string      `json:"ingredientId"`
	Amount       AmountRange `json:"amount"`
	Unit         string      `json:"unit"`
	Prefix       string      `json:"prefix"`
	Suffix       string      `json:"suffix"`
	Text         string      `json:"text"`
	Optional     bool        `json:"optional"`
}

type Ingredient struct {
	RecipeSection   string           `json:"recipeSection"`
	IngredientsList []IngredientData `json:"ingredientsList"`
}

type Serves struct {
	Amount AmountRange `json:"amount"`
	Unit   string      `json:"unit"`
	Text   string      `json:"text"`
}

type Timing struct {
	Qualifier      string      `json:"qualifier"`
	DurationInMins AmountRange `json:"durationInMins"`
	Text           string      `json:"text"`
}

type Instruction struct {
	Description string        `json:"description"`
	Images      []RecipeImage `json:"images"`
}

type CommerceCTA struct {
	SponsorName string `json:"sponsorName"`
	Territory   string `json:"territory"`
	URL         string `json:"url"`
}

type StructuredRecipe struct {
	ID                      string        `json:"id"`
	ComposerID              string        `json:"composerId"`
	CanonicalArticle        string        `json:"canonicalArticle"`
	Title                   string        `json:"title"`
	Description             string        `json:"description"`
	IsAppReady              bool          `json:"isAppReady"`
	FeaturedImage           *RecipeImage  `json:"featuredImage,omitempty"`
	Contributors            []string      `json:"contributors"`
	Byline                  []string      `json:"byline"`
	Ingredients             []Ingredient  `json:"ingredients"`
	SuitableForDietIDs      []string      `json:"suitableForDietIds"`
	CuisineIDs              []string      `json:"cuisineIds"`
	MealTypeIDs             []string      `json:"mealTypeIds"`
	CelebrationIDs          []string      `json:"celebrationIds"`
	UtensilsAndApplianceIDs []string      `json:"utensilsAndApplianceIds"`
	TechniquesUsedIDs       []string      `json:"techniquesUsedIds"`
	DifficultyLevel         string        `json:"difficultyLevel"`
	Serves                  []Serves      `json:"serves"`
	Timings                 []Timing      `json:"timings"`
	Instructions            []Instruction `json:"instructions"`
	CommerceCTAs            []CommerceCTA `json:"commerceCtas,omitempty"`
	BookCredit              string        `json:"bookCredit"`
}

const markdownTemplate = `
# {{.Title}}
By: {{if .Contributors}} {{range .Contributors}}{{.}} {{end}}{{else}}{{range .Byline}}{{.}} {{end}}{{end}}
Originally published at: https://www.theguardian.com/{{.CanonicalArticle}}

{{if .BookCredit}}_{{.BookCredit}}_{{end}}

## Description
{{.Description}}

## Tags
{{if .SuitableForDietIDs}}Meets dietary requirements for: {{range .SuitableForDietIDs}}{{.}},{{end}}{{end}}
{{if .MealTypeIDs}}Suitable meal types: {{range .MealTypeIDs}}{{.}},{{end}}{{end}}
{{if .CelebrationIDs}}Suitable for celebrations: {{range .CelebrationIDs}}{{.}},{{end}}{{end}}
{{if .CuisineIDs}}Cuisine styles: {{range .CuisineIDs}}{{.}},{{end}}{{end}}

## Timings
{{range .Timings}}{{.Qualifier}} {{ .DurationInMins.Render }} minutes
{{end}}

## Ingredients
{{range .Ingredients}}
{{if .RecipeSection }}### {{.RecipeSection}}
{{end}}{{range .IngredientsList}}{{if .Prefix}}{{.Prefix}} {{end}}{{.Amount.Render}} {{.Unit}} {{.Name}}{{if .Suffix}}, {{.Suffix}}{{end}}
{{end}}
{{end}}

## Method
{{range .Instructions}}- {{.Description}}
{{end}}
`

func (r *StructuredRecipe) RenderAsMarkdown() (string, error) {
	tmpl, err := template.New("markdown").Parse(markdownTemplate)
	if err != nil {
		return "", err
	}
	var buf strings.Builder
	err = tmpl.Execute(&buf, r)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
