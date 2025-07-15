package models

import (
	"encoding/json"
	"fmt"
)

type Contributor interface {
	isContributor()
}

type ContributorContributor struct {
	Type       string `json:"type"` // always "contributor"
	TagID      string `json:"tagId"`
	InternalID *int   `json:"internalId,omitempty"`
}

func (c ContributorContributor) isContributor() {}

type ContributorUntagged struct {
	Type string `json:"type"` // always "untagged-contributor"
	Text string `json:"text"`
}

func (c ContributorUntagged) isContributor() {}

type ContributorFreetext struct {
	Type string `json:"type"` // always "freetext"
	Text string `json:"text"`
}

func (c ContributorFreetext) isContributor() {}

/*
*
Because Go doesn't support tagged unions directly, we need a custom unmarshaler to inspect the type field and decode into the correct struct:
*/
func UnmarshalContributor(data []byte) (Contributor, error) {
	var typeDetector struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &typeDetector); err != nil {
		return nil, err
	}

	switch typeDetector.Type {
	case "contributor":
		var c ContributorContributor
		if err := json.Unmarshal(data, &c); err != nil {
			return nil, err
		}
		return c, nil
	case "untagged-contributor":
		var c ContributorUntagged
		if err := json.Unmarshal(data, &c); err != nil {
			return nil, err
		}
		return c, nil
	case "freetext":
		var c ContributorFreetext
		if err := json.Unmarshal(data, &c); err != nil {
			return nil, err
		}
		return c, nil
	default:
		return nil, fmt.Errorf("unknown contributor type: %s", typeDetector.Type)
	}
}

type ContributorList []Contributor

/*
*
Define the ability to unmarshal a slice of Contributor as well, as a method in order to work
seamlessly with json.Unmarshal
*/
func (cl *ContributorList) UnmarshalJSON(data []byte) error {
	var rawItems []json.RawMessage
	if err := json.Unmarshal(data, &rawItems); err != nil {
		return err
	}

	for _, item := range rawItems {
		contrib, err := UnmarshalContributor(item)
		if err != nil {
			return err
		}
		*cl = append(*cl, contrib)
	}
	return nil
}
