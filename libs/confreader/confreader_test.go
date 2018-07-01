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

	[pushbullet]
	token = "dummytoken"
	device = "dummydevice"

	[[series]]
	key = "thunderbirds"
	filenametemplate = "dummyfilenametemplate"

	seriesregexp = ".*"
	seriesreplacement = "Thunderbirds Are Go"
	
	seasonregexp = "^Säsong ([0-9]+)$"
	seasonreplacement = "S$1"

	episoderegexp = "^Avsnitt ([0-9]+)$"
	episodereplacement = "E$1"
`
	_, err := Parse(strings.NewReader(textFile))
	if err != nil {
		t.Error("Did not expect error", err)
	}
}

func TestParseSuperfluesEntries(t *testing.T) {
	textFile := `
	basefolder = "/media/data/Series/"

	unusedKey = "yada"

	[pushbullet]
	token = "dummytoken"
	device = "dummydevice"

	[[series]]
	key = "thunderbirds"
	filenametemplate = "dummyfilenametemplate"

	seriesregexp = ".*"
	seriesreplacement = "Thunderbirds Are Go"
	
	seasonregexp = "^Säsong ([0-9]+)$"
	seasonreplacement = "S$1"

	episoderegexp = "^Avsnitt ([0-9]+)$"
	episodereplacement = "E$1"
`
	_, err := Parse(strings.NewReader(textFile))
	if err == nil {
		t.Error("Did expect error")
	}
}
