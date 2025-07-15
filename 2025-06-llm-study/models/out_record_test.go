package models

import (
	"io"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestParseLLMContent(t *testing.T) {
	f, err := os.Open("fixture/response.md")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fixtureData, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	fixtureString := string(fixtureData)
	result := ParseLLMContent(&fixtureString)
	spew.Dump(result)

	if result.MarkdownStart != 12 {
		t.Errorf("Expected 12 for MarkdownStart, got %d", result.MarkdownStart)
	}
	if result.MarkdownEnd != 1513 {
		t.Errorf("Expected 1513 for MarkdownEnd, got %d", result.MarkdownEnd)
	}
	if result.NotesStart != 1517 {
		t.Errorf("Expected 1517 for NotesStart, got %d", result.NotesStart)
	}
	if result.NotesEnd != 3056 {
		t.Errorf("Expected 3056 for NotesEnd, got %d", result.NotesEnd)
	}
}

func TestProcessResultStruct(t *testing.T) {
	f, err := os.Open("fixture/response.md")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fixtureData, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	fixtureString := string(fixtureData)

	rec := Record{}
	rec.ProcessResultText(&fixtureString)

	if rec.AnnotatedText != "# Chicken drumsticks with pomegranate and oregano\nBy: <!HEALTH Missing author>\nOriginally published at: https://www.theguardian.com/global/2016/apr/18/easy-ottolenghi-meat-recipes\n\n## Description\n<!HEALTH Missing description>\n\n## Tags\n\nSuitable meal types: main-course,\n\n## Timings\n<!HEALTH Missing timings>\n\n## Ingredients\n\n2 tbsp olive oil\n8  chicken drumsticks\n  Flaky sea salt\n500 g banana shallots\n1 head garlic\n2 tbsp pomegranate molasses\n3 tbsp light soy sauce\n1 tbsp maple syrup\n 4 cm piece fresh ginger\n10 sprigs oregano\n0.50  small pomegranate\n\n## Method\n- Heat the oven to 200C/390F/gas mark 6.\n- Heat the oil in a large frying pan on a medium-high flame, then lay the drumsticks in the pan, sprinkle over half a teaspoon of salt and fry, turning regularly, for 10 minutes, until golden-brown all over.\n- Transfer to a large bowl and leave the pan on the heat.\n- Fry the shallots for four minutes, shaking the pan a few times.\n- Add the garlic, fry for another minute, until golden-brown, add to the chicken bowl, then combine with the molasses, soy sauce, maple syrup, ginger, oregano sprigs and 50ml water.\n- Pour the contents of the bowl into a 26cm x 36cm oven dish or tray and cover tightly with aluminium foil.\n- Roast for 20 minutes, then take off the foil, stir everything together and roast for 10 minutes more, until the chicken is cooked through and the shallots and garlic are soft.\n- Remove from the oven, stir in the chopped oregano leaves and pomegranate seeds, and serve." {
		spew.Dump(rec)
		t.Error("AnnotatedText was not properly extracted")
	}
	if rec.ModelNotes != "The following issues were found:\n\n1. Missing author\n2. Missing description\n3. Missing timings\n4. Ingredient 'Flaky sea salt' does not specify a quantity, but it is a seasoning so it could be acceptable. However, for consistency, it would be better to specify 'to taste' or a suggested quantity.\n5. Ingredient '1 head garlic' could be clearer by specifying cloves, for example '1 head of garlic, cloves separated'.\n6. Ingredient '0.50 small pomegranate' should not have a decimal point as it's not a standard unit of measurement. It should be '½ small pomegranate'.\n7. The method step 'Fry the shallots for four minutes, shaking the pan a few times.' could be clearer by specifying the oil or fat used for frying.\n8. The method step 'Add the garlic, fry for another minute, until golden-brown, add to the chicken bowl, then combine with the molasses, soy sauce, maple syrup, ginger, oregano sprigs and 50ml water.' is too long and combines multiple steps. It should be broken down into smaller steps.\n9. The method step 'Pour the contents of the bowl into a 26cm x 36cm oven dish or tray and cover tightly with aluminium foil.' does not specify the oven temperature or cooking time, although it mentions roasting in the next step. It would be clearer to include this information with this step.\n10. The method step 'Remove from the oven, stir in the chopped oregano leaves and pomegranate seeds, and serve.' mentions 'chopped oregano leaves', but the ingredients list only includes 'oregano sprigs'. This inconsistency should be resolved." {
		spew.Dump(rec)
		t.Error("ModelNotes was not properly extracted")
	}
	if len(rec.AnnotationSummary) != int(rec.AnnotationCount) {
		t.Errorf("AnnotationCount was not equal to the length of AnnotationSummary, %d vs %d", rec.AnnotationCount, len(rec.AnnotationSummary))
	}
	if rec.AnnotationSummary[0] != "Missing author" {
		t.Error("Missing first expected annotation")
	}
	if rec.AnnotationSummary[1] != "Missing description" {
		t.Error("Missing first expected annotation")
	}
	if rec.AnnotationSummary[2] != "Missing timings" {
		t.Error("Missing first expected annotation")
	}
}
