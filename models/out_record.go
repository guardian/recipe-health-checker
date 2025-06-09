package models

import (
	"bufio"
	"regexp"
	"strings"
	"time"
)

type ParsedReturnContent struct {
	MarkdownStart int //character position of the content
	MarkdownEnd   int
	NotesStart    int
	NotesEnd      int
}

func ParseLLMContent(content *string) *ParsedReturnContent {
	out := ParsedReturnContent{}

	r := strings.NewReader(*content)
	scanner := bufio.NewScanner(r)

	inMarkdown := false
	ctr := 0
	for scanner.Scan() {
		line := scanner.Text()
		if !inMarkdown && strings.HasPrefix(line, "```markdown") {
			out.MarkdownStart = ctr + len(line) + 1 //markdown content starts from the next line (include the newline char)
			inMarkdown = true
		} else if inMarkdown && strings.HasPrefix(line, "```") {
			out.MarkdownEnd = ctr //markdown content ends at this line
			out.NotesStart = ctr + len(line) + 1
		}
		ctr += len(line) + 1
	}
	out.NotesEnd = ctr - 1
	return &out
}

type Record struct {
	OriginalRecipeSnapshot string    `json:"snapshot"`
	AnnotatedText          string    `json:"annotated_text"`
	ModelNotes             string    `json:"model_notes"`
	AnnotationCount        uint32    `json:"annotation_count"`
	AnnotationSummary      []string  `json:"annotation_summary"`
	RecipeId               string    `json:"recipe_id"`
	ComposerId             string    `json:"composer_id"`
	Timestamp              time.Time `json:"timestamp"`
	InputTokensUsed        int       `json:"input_tokens_used"`
	OutputTokensUsed       int       `json:"output_tokens_used"`
	ModelUsed              string    `json:"model_used"`
}

/*
*
ProcessResultText will parse the content that came from the ML model, separating
notes and counting annotations
*/
func (r *Record) ProcessResultText(resultText *string) {
	locations := ParseLLMContent(resultText)
	r.AnnotatedText = strings.TrimSpace((*resultText)[locations.MarkdownStart:locations.MarkdownEnd])
	r.ModelNotes = strings.TrimSpace((*resultText)[locations.NotesStart:locations.NotesEnd])

	notextractor := regexp.MustCompile(`(?m)<!HEALTH ([^>]+)>`)

	notes := notextractor.FindAllStringSubmatch(r.AnnotatedText, -1)
	if notes != nil {
		r.AnnotationCount = uint32(len(notes))
		r.AnnotationSummary = make([]string, r.AnnotationCount)
		for i, match := range notes {
			r.AnnotationSummary[i] = match[1]
		}
	}
}
