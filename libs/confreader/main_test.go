package confreader

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	textFile := `
`
	_, err := Parse(strings.NewReader(textFile))
	if err == nil {
		t.Error("Empty config should have generated error")
	}
}

func TestParseOk(t *testing.T) {
	textFile := `
	basefolder = "/media/data/Series/"

	[[series]]
	key = "thunderbirds"

	[series.series]
	regexp = ".*"
	replacement = "Thunderbirds Are Go"
	
	[series.season]
	regexp = "^Säsong ([0-9]+)$"
	replacement = "S$1"

	[series.episode]
	regexp = "^Avsnitt ([0-9]+)$"
	replacement = "E$1"
`
	_, err := Parse(strings.NewReader(textFile))
	if err != nil {
		t.Error("Did not expect error", err)
	}
}
