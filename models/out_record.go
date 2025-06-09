package models

import (
	"bufio"
	"github.com/davecgh/go-spew/spew"
	"log"
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
		if !inMarkdown && (strings.HasPrefix(line, "```markdown") || strings.HasPrefix(line, "```json")) {
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
	AnnotationSection      []string  `json:"annotation_section"`
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
	textLen := len(*resultText)
	if textLen == 0 {
		log.Printf("WARNING: cannot process result as it was empty")
	}

	locations := ParseLLMContent(resultText)
	spew.Dump(locations)
	r.AnnotatedText = strings.TrimSpace((*resultText)[locations.MarkdownStart:locations.MarkdownEnd])
	// Guard against inaccurate calculations pushing is over
	notesEnd := locations.NotesEnd
	if notesEnd > textLen {
		notesEnd = textLen - 2
	}
	if locations.NotesStart < notesEnd && locations.NotesEnd > 0 {
		r.ModelNotes = strings.TrimSpace((*resultText)[locations.NotesStart:notesEnd])
	}

	notextractor := regexp.MustCompile(`(?m)<!HEALTH(:\S+)* ([^>]+)>`)

	notes := notextractor.FindAllStringSubmatch(r.AnnotatedText, -1)
	if notes != nil {
		r.AnnotationCount = uint32(len(notes))
		r.AnnotationSummary = make([]string, r.AnnotationCount)
		r.AnnotationSection = make([]string, r.AnnotationCount)
		for i, match := range notes {
			if len(match) == 2 {
				r.AnnotationSummary[i] = match[1]
			} else if len(match) == 3 {
				r.AnnotationSummary[i] = match[2]
				r.AnnotationSection[i] = match[1]
			} else if len(match) > 3 {
				log.Printf("WARNING Got more matches than expected on health tag %d", i)
			} else {
				log.Printf("WARNING Got health tag %d without any content", i)
			}
		}
	}
}
