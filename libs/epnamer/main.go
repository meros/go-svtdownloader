package epnamer

import (
	"bytes"
	"regexp"

	"github.com/alecthomas/template"
	"github.com/meros/go-svtdownloader/libs/eplister"
)

type Replacement struct {
	Re          *regexp.Regexp
	Replacement string
}

// Options contain options for how this series should be translated into a filename
type Options struct {
	// Series is used to translate season string into something else if non-nil
	Series *Replacement
	// Season is used to translate season string into something else if non-nil
	Season *Replacement
	// Season is used to translate season string into something else if non-nil
	Episode *Replacement
	// Template is the file output template
	// An example would be "{{.Series}}/{{.Series}} - {{.Season}}{{.Episode}}"
	Template *template.Template
}

// Filename translate an episode to a filename
func Filename(ep eplister.Episode, options Options) (filename string, err error) {
	series := ep.Series
	if options.Series != nil {
		series = options.Series.Re.ReplaceAllString(series, options.Series.Replacement)
	}

	season := ep.Season
	if options.Season != nil {
		season = options.Season.Re.ReplaceAllString(season, options.Season.Replacement)
	}

	episode := ep.Episode
	if options.Episode != nil {
		episode = options.Episode.Re.ReplaceAllString(episode, options.Episode.Replacement)
	}

	outfilename := bytes.NewBufferString("")
	err = options.Template.Execute(outfilename, &struct {
		Series  string
		Season  string
		Episode string
	}{series, season, episode})
	if err != nil {
		return "", err
	}

	return outfilename.String(), nil
}
