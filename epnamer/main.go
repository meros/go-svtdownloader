package epnamer

import (
	"bytes"

	"github.com/meros/go-svtdownloader/confreader"
	"github.com/meros/go-svtdownloader/eplister"
	"text/template"
)

// Options contain options for how this series should be translated into a filename
type Options struct {
	Series  confreader.Transformer
	Season  confreader.Transformer
	Episode confreader.Transformer
	// Template is the file output template
	// An example would be "{{.Series}}/{{.Series}} - {{.Season}}{{.Episode}}"
	Template *template.Template
}

// Filename translate an episode to a filename
func Filename(ep eplister.Episode, options Options) (filename string, err error) {
	series := options.Series.Transform(ep.Series)
	season := options.Season.Transform(ep.Season)
	episode := options.Episode.Transform(ep.Episode)

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
