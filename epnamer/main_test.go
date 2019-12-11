package epnamer

import (
	"github.com/meros/go-svtdownloader/eplister"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"text/template"
)

type identityTransformer struct{}

func (ret *identityTransformer) Transform(text string) string {
	return text
}

func TestFilenameSimple(t *testing.T) {
	url, err := url.Parse("httsp://dummy.se")
	if err != nil {
		t.Error("Failed to create dummy url", err)
		return
	}

	tmplt, err := template.New("Test template").Parse("{{.Series}} {{.Season}} {{.Episode}}")
	if err != nil {
		t.Error("Template failed with err", err)
		return
	}

	filename, err := Filename(eplister.Episode{
		"Thunderbirds",
		"Säsong 12",
		"Avsnitt 14",
		*url},
		Options{
			Series:   &identityTransformer{},
			Season:   &identityTransformer{},
			Episode:  &identityTransformer{},
			Template: tmplt})

	if err != nil {
		t.Error("Filename failed with err", err)
		return
	}

	if strings.Compare(filename, "Thunderbirds Säsong 12 Avsnitt 14") != 0 {
		t.Error("Not expected result", filename)
	}
}

type regexpTransformer struct {
	re *regexp.Regexp
	rs string
}

func (ret *regexpTransformer) Transform(text string) string {
	return ret.re.ReplaceAllString(text, ret.rs)
}

func TestFilenameReplacements(t *testing.T) {
	url, err := url.Parse("httsp://dummy.se")
	if err != nil {
		t.Error("Failed to create dummy url", err)
		return
	}

	tmplt, err := template.New("Test template").Parse("{{.Series}}/{{.Series}} - {{.Season}}{{.Episode}}")
	if err != nil {
		t.Error("Template failed with err", err)
		return
	}

	filename, err := Filename(eplister.Episode{
		"Thunderbirds",
		"Säsong 12",
		"Avsnitt 14",
		*url},
		Options{
			Series: &regexpTransformer{
				regexp.MustCompile("^Thunderbirds$"),
				"Thunderbirds Are Go"},
			Season: &regexpTransformer{
				regexp.MustCompile("^Säsong ([0-9]+)$"),
				"S$1"},
			Episode: &regexpTransformer{
				regexp.MustCompile("^Avsnitt ([0-9]+)$"),
				"E$1"},
			Template: tmplt})

	if err != nil {
		t.Error("Filename failed with err", err)
		return
	}

	if strings.Compare(filename, "Thunderbirds Are Go/Thunderbirds Are Go - S12E14") != 0 {
		t.Error("Not expected result", filename)
	}
}
