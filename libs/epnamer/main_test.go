package epnamer

import (
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/meros/go-svtdownloader/libs/eplister"
)

func TestFilenameSimple(t *testing.T) {
	url, err := url.Parse("httsp://dummy.se")
	if err != nil {
		t.Error("Failed to create dummy url", err)
		return
	}

	filename, err := Filename(eplister.Episode{
		"Thunderbirds",
		"Säsong 12",
		"Avsnitt 14",
		*url},
		Options{
			nil,
			nil,
			nil,
			"{{.Series}} {{.Season}} {{.Episode}}"})

	if err != nil {
		t.Error("Filename failed with err", err)
		return
	}

	if strings.Compare(filename, "Thunderbirds Säsong 12 Avsnitt 14") != 0 {
		t.Error("Not expected result", filename)
	}
}

func TestFilenameReplacements(t *testing.T) {
	url, err := url.Parse("httsp://dummy.se")
	if err != nil {
		t.Error("Failed to create dummy url", err)
		return
	}

	filename, err := Filename(eplister.Episode{
		"Thunderbirds",
		"Säsong 12",
		"Avsnitt 14",
		*url},
		Options{
			Series: &Replacement{
				*regexp.MustCompile("^Thunderbirds$"),
				"Thunderbirds Are Go"},
			Season: &Replacement{
				*regexp.MustCompile("^Säsong ([0-9]+)$"),
				"S$1"},
			Episode: &Replacement{
				*regexp.MustCompile("^Avsnitt ([0-9]+)$"),
				"E$1"},
			TemplateString: "{{.Series}}/{{.Series}} - {{.Season}}{{.Episode}}"})

	if err != nil {
		t.Error("Filename failed with err", err)
		return
	}

	if strings.Compare(filename, "Thunderbirds Are Go/Thunderbirds Are Go - S12E14") != 0 {
		t.Error("Not expected result", filename)
	}
}
