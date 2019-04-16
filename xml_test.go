package suggest

import (
	"encoding/xml"
	"strings"
	"testing"
)

var (
	testXMLDATA = `<?xml version="1.0"?><toplevel><CompleteSuggestion><suggestion data="golang"/></CompleteSuggestion><CompleteSuggestion><suggestion data="golang tutorial"/></CompleteSuggestion><CompleteSuggestion><suggestion data="golang json"/></CompleteSuggestion><CompleteSuggestion><suggestion data="golang interface"/></CompleteSuggestion><CompleteSuggestion><suggestion data="golang time"/></CompleteSuggestion><CompleteSuggestion><suggestion data="golang playground"/></CompleteSuggestion><CompleteSuggestion><suggestion data="golang map"/></CompleteSuggestion><CompleteSuggestion><suggestion data="golang download"/></CompleteSuggestion><CompleteSuggestion><suggestion data="golang switch"/></CompleteSuggestion><CompleteSuggestion><suggestion data="golang for"/></CompleteSuggestion></toplevel>`
)

func TestDecode(t *testing.T) {
	dec := xml.NewDecoder(strings.NewReader(testXMLDATA))

	var s GoogleSuggestion
	err := dec.Decode(&s.TopLevel)
	if err != nil {
		t.Errorf("\nerror: %v\n", err)
		return
	}

	l := len(s.TopLevel.CompleteSuggestion)
	if l != 10 {
		t.Errorf("\ngot : %d, want: %d\n", l, 10)
		return
	}

	data := []string{
		"golang",
		"golang tutorial",
		"golang json",
		"golang interface",
		"golang time",
		"golang playground",
		"golang map",
		"golang download",
		"golang switch",
		"golang for",
	}

	for i, v := range s.TopLevel.CompleteSuggestion {
		if v.Suggestion.Data != data[i] {
			t.Errorf("\ngot : %v, want: %v\n", v.Suggestion.Data, data[i])
		}
	}
}
