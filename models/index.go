package models

import "fmt"

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
