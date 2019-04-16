package suggest

import (
	"encoding/xml"
	"io"
)

type Suggestion struct {
	Data string `xml:"data,attr" json:"data"`
}

type CompleteSuggestion struct {
	Suggestion Suggestion `xml:"suggestion" json:"suggestion"`
}

type TopLevel struct {
	CompleteSuggestion []CompleteSuggestion `xml:"CompleteSuggestion" json:"CompleteSuggestion"`
}

type GoogleSuggestion struct {
	TopLevel TopLevel `xml:"toplevel" json:"toplevel"`
}

// XMLDecode XML convert to structure
func XMLDecode(r io.Reader) (GoogleSuggestion, error) {
	var gs GoogleSuggestion
	dec := xml.NewDecoder(r)
	err := dec.Decode(&gs.TopLevel)
	return gs, err
}
