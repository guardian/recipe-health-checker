package models

import (
	"fmt"
	"math/rand"
	"time"
)

type IndexEntry struct {
	Checksum         string `json:"checksum"`
	RecipeID         string `json:"recipeUID"`
	CapiID           string `json:"capiArticleId"`
	SponsorshipCount int    `json:"sponsorshipCount"`
}

func (i *IndexEntry) GetUrl(baseUrl string) string {
	return fmt.Sprintf("%s/content/%s", baseUrl, i.Checksum)
}

type RecipeIndex struct {
	SchemaVersion int           `json:"schemaVersion"`
	Recipes       []*IndexEntry `json:"recipes"`
}

func (index *RecipeIndex) RandomisedSample(sampleCount int) []*IndexEntry {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	maxValue := len(index.Recipes)
	output := make([]*IndexEntry, 0)

	for range sampleCount {
		idx := int32((r.Float64() * float64(maxValue)))
		output = append(output, index.Recipes[idx])
	}
	return output
}
