package llm

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/guardian/recipe-health-checker/models"
)

type LLM struct {
	client             *bedrockruntime.Client
	ModelId            string
	systemContentBlock types.SystemContentBlock
}

type RecipeFormat int

const (
	Json     RecipeFormat = iota
	Markdown RecipeFormat = iota
)

func New(ctx context.Context, modelName string, region string) *LLM {
	system := types.SystemContentBlockMemberText{
		Value: `Your job is to proof-read recipes for publication in an app.  You will be presented with a recipe in Markdown format, containing clearly deliniated sections for description, recipe tags, recipe timings, required ingredients (possibly broken down by section) and method steps.
		
		You should return the recipe data you were given and flag any problems in the text. If you find any examples of the kind of problems listed below you should insert a marker into the text using this format: <!HEALTH:{recipe-section} A description of the issue>.

		Each recipe should, at the very least, have a title at the top; at least one author following 'By:'; at least one kind of tag (the more, the better), a full method and every ingredient mentioned in the method should be included in the ingredients list.
		If vital information like title, author or description are missing then these must be flagged.

		Please pay special attention to the ingredients list.  Ingredients are rendered in the format {quantity} {unit} {name}, {optional suffix}. If you see any ingredients not following this pattern please flag this is a problem.  Extra whitespace should be ignored.
		For example:
		- '4 tbsp white vermouth, such as Noilly Prat' is a good ingredient line because it clearly contains the quantity, unit, name and some extra information at the end.
		- '(10g) 6 tbsp fresh oregano leaves' is also a good ingredient line because the quantity has been specified twice in different units (10 grammes and 6 tbsp)
		- '  Salt and black pepper, to taste' is also a good ingredient line because salt and pepper are added to many recipes as seasoning without needing a defined quantity.  Be careful, this exception only applies a few ingredients (like salt, pepper, herbs like chopped coriander, etc.)
		- '  Freshly ground back pepper' is also fine because it is used as seasoning and does not need a defined quantity
		- 'finely 10 g chopped fresh flat-leaf parsley' is a bad ingredient line because it should read '10g flat-leaf parsley, fresh and finely chopped'
		- 'long 2  red chillies, finely chopped' is a bad ingredient line because it should read '2 long red chillies, finely chopped'
		
		Over-long or over-short method steps should also be flagged.  A method step should contain only one specific operation for cooking the recipe, for example:
		- 'Occasionally scrape down the mixture from the sides of the bowl with a rubber spatula.' is a good length of a step
		- 'Break the egg into a small bowl and briefly blend with a fork. Add the egg to the butter and sugar, a little at a time, beating continuously. (Should the mixture curdle, add a spoonful of the flour.)' is also a good length
		- 'Add the flour and baking powder, turning the mixture slowly. Stir the apricots and water into the mixture. (This will alter the consistency alarmingly, but do not worry, all will come good in the oven.) Transfer the batter to the baking dish and bake for 35 minutes.' is not a good step, because it combines three seperate steps into one sentence making it too long

		Do not flag capitalization in tag fields as these are re-rendered by the app

		If no issues are found, do not return the whole document, simply return <!HEALTH No issues found :+1>
		If you determine any text to be in violation of content filtering please remove it from the output and indicate this using <!HEALTH redacted>
		`,
	}

	return &LLM{
		client:             newClient(ctx, region),
		ModelId:            modelName,
		systemContentBlock: &system,
	}
}

func newClient(ctx context.Context, region string) *bedrockruntime.Client {
	sdkConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		panic(err)
	}
	brClient := bedrockruntime.NewFromConfig(sdkConfig)
	return brClient
}

func (m *LLM) generateInput(content string, format RecipeFormat) *bedrockruntime.ConverseInput {
	sanitiser := regexp.MustCompile("```")
	sanitisedContent := sanitiser.ReplaceAllString(content, "") //make sure that markdown can't escape fencing

	contentFmt := "unknown"
	switch format {
	case Json:
		contentFmt = "JSON"
		break
	case Markdown:
		contentFmt = "Markdown"
		break
	}

	requestMsg := types.Message{
		Role: types.ConversationRoleUser,
		Content: []types.ContentBlock{
			&types.ContentBlockMemberText{
				Value: fmt.Sprintf(
					"Here is a recipe in %s format.  Please check it over and show me where there may be problems.\n\n```%s\n%s\n```",
					contentFmt,
					strings.ToLower(contentFmt),
					sanitisedContent,
				)},
		},
	}

	return &bedrockruntime.ConverseInput{
		ModelId: &m.ModelId,
		System: []types.SystemContentBlock{
			m.systemContentBlock,
			//&types.SystemContentBlockMemberCachePoint{Value: types.CachePointBlock{Type: types.CachePointTypeDefault}},
		},
		Messages: []types.Message{
			requestMsg,
		},
	}
}

func extractTextBlocks(blocks *[]types.ContentBlock) []types.ContentBlockMemberText {
	output := make([]types.ContentBlockMemberText, 0)

	for _, b := range *blocks {
		if textBlock, isText := b.(*types.ContentBlockMemberText); isText {
			output = append(output, *textBlock)
		}
	}
	return output
}

func extractText(blocks *[]types.ContentBlock) string {
	var out strings.Builder

	textBlocks := extractTextBlocks(blocks)
	for _, b := range textBlocks {
		out.Write([]byte(b.Value))
	}
	return out.String()
}

func (m *LLM) RequestReview(ctx context.Context, content string, format RecipeFormat) (*models.Record, error) {
	req := m.generateInput(content, format)

	response, err := m.client.Converse(ctx, req)
	if err != nil {
		return nil, err
	}

	if msg, isMsg := response.Output.(*types.ConverseOutputMemberMessage); isMsg {
		text := extractText(&msg.Value.Content)

		result := &models.Record{
			OriginalRecipeSnapshot: content,
			Timestamp:              time.Now(),
			InputTokensUsed:        int(*response.Usage.InputTokens),
			OutputTokensUsed:       int(*response.Usage.OutputTokens),
		}
		result.ProcessResultText(&text)

		if response.StopReason != types.StopReasonEndTurn {
			return result, fmt.Errorf("The model stopped unexpectedly: %s", response.StopReason)
		}
		return result, nil
	} else {
		return nil, errors.New("No text was returned from the model")
	}
}
