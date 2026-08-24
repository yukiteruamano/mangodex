package mangodex

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalisedStrings_Unmarshal_Map(t *testing.T) {
	j := `{"en":"Hello","ja":"こんにちは"}`
	var ls LocalisedStrings
	require.NoError(t, json.Unmarshal([]byte(j), &ls))
	assert.Equal(t, "Hello", ls.GetLocalString("en"))
	assert.Equal(t, "こんにちは", ls.GetLocalString("ja"))
}

func TestLocalisedStrings_Unmarshal_ArrayOfMaps(t *testing.T) {
	j := `[{"en":"Hello"},{"ja":"こんにちは"}]`
	var ls LocalisedStrings
	require.NoError(t, json.Unmarshal([]byte(j), &ls))
	assert.Equal(t, "Hello", ls.Values["en"])
	assert.Equal(t, "こんにちは", ls.Values["ja"])
}

func TestLocalisedStrings_Unmarshal_Null(t *testing.T) {
	var ls LocalisedStrings
	require.NoError(t, json.Unmarshal([]byte(`null`), &ls))
	assert.Empty(t, ls.Values)
}

func TestLocalisedStrings_GetLocalString_Fallback(t *testing.T) {
	ls := LocalisedStrings{Values: map[string]string{"en": "Hello"}}
	assert.Equal(t, "Hello", ls.GetLocalString("fr"))
	ls2 := LocalisedStrings{Values: map[string]string{}}
	assert.Equal(t, "", ls2.GetLocalString("en"))
}

func TestLocalisedStrings_Marshal(t *testing.T) {
	ls := LocalisedStrings{Values: map[string]string{"en": "Hi"}}
	b, err := json.Marshal(ls)
	require.NoError(t, err)
	assert.Contains(t, string(b), `"en":"Hi"`)

	ls2 := LocalisedStrings{Values: nil}
	b2, err := json.Marshal(ls2)
	require.NoError(t, err)
	assert.Equal(t, `{}`, string(b2))
}

func TestRelationship_Unmarshal_Manga(t *testing.T) {
	j := `{"id":"manga-id","type":"manga","attributes":{"title":{"en":"Title"},"altTitles":{},"description":{},"isLocked":false,"links":{},"originalLanguage":"ja","state":"published","version":1,"createdAt":"2020-01-01T00:00:00","updatedAt":"2020-01-01T00:00:00","tags":[]}}`
	var rel Relationship
	require.NoError(t, json.Unmarshal([]byte(j), &rel))
	assert.Equal(t, "manga", rel.Type)
	_, ok := rel.Attributes.(*MangaAttributes)
	assert.True(t, ok)
}

func TestRelationship_Unmarshal_Author(t *testing.T) {
	j := `{"id":"auth-id","type":"author","attributes":{"name":"Author","imageUrl":null,"biography":{},"version":1,"createdAt":"2020-01-01T00:00:00","updatedAt":"2020-01-01T00:00:00"}}`
	var rel Relationship
	require.NoError(t, json.Unmarshal([]byte(j), &rel))
	assert.Equal(t, "author", rel.Type)
	_, ok := rel.Attributes.(*AuthorAttributes)
	assert.True(t, ok)
}

func TestRelationship_Unmarshal_Unknown(t *testing.T) {
	j := `{"id":"x","type":"unknown_type","attributes":{"foo":"bar"}}`
	var rel Relationship
	require.NoError(t, json.Unmarshal([]byte(j), &rel))
	_, ok := rel.Attributes.(*json.RawMessage)
	assert.True(t, ok)
}

func TestRelationship_Unmarshal_NullAttributes(t *testing.T) {
	j := `{"id":"x","type":"manga","attributes":null}`
	var rel Relationship
	require.NoError(t, json.Unmarshal([]byte(j), &rel))
	assert.Equal(t, "x", rel.ID)
}

func TestErrorResponse_GetErrors(t *testing.T) {
	er := ErrorResponse{Result: "error", Errors: []Error{{Title: "T", Detail: "D"}, {Title: "T2", Detail: "D2"}}}
	s := er.GetErrors()
	assert.Contains(t, s, "T: D")
	assert.Contains(t, s, "T2: D2")
}

func TestTag_GetName(t *testing.T) {
	tag := Tag{Attributes: TagAttributes{Name: LocalisedStrings{Values: map[string]string{"en": "Action"}}}}
	assert.Equal(t, "Action", tag.GetName("en"))
}

func TestResponse_GetResult(t *testing.T) {
	r := Response{Result: "ok"}
	assert.Equal(t, "ok", r.GetResult())
	er := ErrorResponse{Result: "error"}
	assert.Equal(t, "error", er.GetResult())
}

func TestManga_GetTitle_Fallback(t *testing.T) {
	m := Manga{Attributes: MangaAttributes{
		Title:     LocalisedStrings{Values: map[string]string{}},
		AltTitles: LocalisedStrings{Values: map[string]string{"en": "Alt"}},
	}}
	assert.Equal(t, "Alt", m.GetTitle("en"))
	m2 := Manga{Attributes: MangaAttributes{Title: LocalisedStrings{Values: map[string]string{"en": "Main"}}}}
	assert.Equal(t, "Main", m2.GetTitle("en"))
}

func TestChapter_GetChapterNum(t *testing.T) {
	ch := Chapter{Attributes: ChapterAttributes{Chapter: strPtr("10")}} //nolint:modernize
	assert.Equal(t, "10", ch.GetChapterNum())
	ch2 := Chapter{Attributes: ChapterAttributes{}}
	assert.Equal(t, "-", ch2.GetChapterNum())
}

func strPtr(s string) *string { return &s } //nolint:modernize
