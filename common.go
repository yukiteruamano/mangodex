package mangodex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"strings"
)

// ResponseType is an interface for API responses.
type ResponseType interface {
	GetResult() string
}

// Response is a plain response containing only Result.
type Response struct {
	Result string `json:"result"`
}

func (r *Response) GetResult() string {
	return r.Result
}

// Relationship contains relationships with optional expanded attributes.
type Relationship struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes any    `json:"attributes"`
}

func (a *Relationship) UnmarshalJSON(data []byte) error {
	typ := struct {
		ID         string          `json:"id"`
		Type       string          `json:"type"`
		Attributes json.RawMessage `json:"attributes"`
	}{}
	if err := json.Unmarshal(data, &typ); err != nil {
		return err
	}

	var err error
	switch typ.Type {
	case MangaRel:
		a.Attributes = &MangaAttributes{}
	case AuthorRel, ArtistRel:
		a.Attributes = &AuthorAttributes{}
	case ScanlationGroupRel:
		a.Attributes = &ScanlationGroupAttributes{}
	case CoverArtRel:
		a.Attributes = &CoverAttributes{}
	case TagRel:
		a.Attributes = &TagAttributes{}
	case UserRel:
		a.Attributes = &UserAttributes{}
	case ChapterRel:
		a.Attributes = &ChapterAttributes{}
	case CustomListRel:
		a.Attributes = &CustomListAttributes{}
	default:
		a.Attributes = &json.RawMessage{}
	}

	a.ID = typ.ID
	a.Type = typ.Type
	if len(typ.Attributes) > 0 && !bytes.Equal(typ.Attributes, []byte("null")) {
		if err = json.Unmarshal(typ.Attributes, a.Attributes); err != nil {
			return fmt.Errorf("error unmarshalling relationship of type %s: %w: %s", typ.Type, err, string(data))
		}
	}
	return err
}

// LocalisedStrings wraps a map of language code to string.
type LocalisedStrings struct {
	Values map[string]string
}

func (l *LocalisedStrings) UnmarshalJSON(data []byte) error {
	l.Values = map[string]string{}
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		return nil
	}
	if err := json.Unmarshal(data, &l.Values); err == nil {
		return nil
	}
	var locals []map[string]string
	if err := json.Unmarshal(data, &locals); err != nil {
		return fmt.Errorf("error unmarshalling localisedstring: %w", err)
	}
	for _, entry := range locals {
		maps.Copy(l.Values, entry)
	}
	return nil
}

func (l LocalisedStrings) MarshalJSON() ([]byte, error) {
	if l.Values == nil {
		return json.Marshal(map[string]string{})
	}
	return json.Marshal(l.Values)
}

// GetLocalString returns the string for langCode or the first available.
func (l *LocalisedStrings) GetLocalString(langCode string) string {
	if s, ok := l.Values[langCode]; ok {
		return s
	}
	for _, v := range l.Values {
		return v
	}
	return ""
}

// Tag contains information on a tag.
type Tag struct {
	ID            string         `json:"id"`
	Type          string         `json:"type"`
	Attributes    TagAttributes  `json:"attributes"`
	Relationships []Relationship `json:"relationships"`
}

// GetName returns the tag name for langCode.
func (t *Tag) GetName(langCode string) string {
	return t.Attributes.Name.GetLocalString(langCode)
}

// TagAttributes holds attributes for a Tag.
type TagAttributes struct {
	Name        LocalisedStrings `json:"name"`
	Description LocalisedStrings `json:"description"`
	Group       string           `json:"group"`
	Version     int              `json:"version"`
}

// ErrorResponse is the typical error response.
type ErrorResponse struct {
	Result string  `json:"result"`
	Errors []Error `json:"errors"`
}

func (er *ErrorResponse) GetResult() string {
	return er.Result
}

// GetErrors returns a formatted string of all errors.
// perf: uses strings.Builder with pre-grow to reduce allocs
func (er *ErrorResponse) GetErrors() string {
	if len(er.Errors) == 0 {
		return ""
	}
	var b strings.Builder
	// Estimate ~32 bytes per error
	b.Grow(len(er.Errors) * 32)
	for _, e := range er.Errors {
		b.WriteString(e.Title)
		b.WriteString(": ")
		b.WriteString(e.Detail)
		b.WriteByte('\n')
	}
	return b.String()
}

// Error holds details of an error.
type Error struct {
	ID     string `json:"id"`
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// ListResponse is a generic paginated list response.
type ListResponse[T any] struct {
	Result   string `json:"result"`
	Response string `json:"response"`
	Data     []T    `json:"data"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
	Total    int    `json:"total"`
}

func (r *ListResponse[T]) GetResult() string { return r.Result }

// SingleResponse is a generic single entity response.
type SingleResponse[T any] struct {
	Result   string `json:"result"`
	Response string `json:"response"`
	Data     T      `json:"data"`
}

func (r *SingleResponse[T]) GetResult() string { return r.Result }
