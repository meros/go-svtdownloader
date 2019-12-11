package confreader

import (
	"fmt"
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
	filenametemplate = "dummyfilenametemplate"

	seriesregexp = ".*"
	seriesreplacement = "Thunderbirds Are Go"
	
	seasonregexp = "^Säsong ([0-9]+)$"
	seasonreplacement = "S$1"

	episoderegexp = "^Avsnitt ([0-9]+)$"
	episodereplacement = "E$1"
`
	config, err := Parse(strings.NewReader(textFile))
	if err != nil {
		t.Error("Did not expect error", err)
		return
	}

	fmt.Println(config)

	if len(config.Series) != 1 {
		t.Error("Did not expect empty series")
		return
	}

	if strings.Compare(config.Series[0].SeasonTransformer.Transform("Säsong 42"), "S42") != 0 {
		t.Error("SeasonTransformer did not produce expected output")
		return
	}

	if strings.Compare(config.Series[0].EpisodeTransformer.Transform("Avsnitt 42"), "E42") != 0 {
		t.Error("SeasonTransformer did not produce expected output")
		return
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
