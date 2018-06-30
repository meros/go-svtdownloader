package confreader

import (
	"errors"
	"io"
	"regexp"

	"github.com/BurntSushi/toml"
)

type Replacement struct {
	Regex       string
	Replacement string
}

type Series struct {
	Key              string
	FilenameTemplate string
	Series           *Replacement
	Season           *Replacement
	Episode          *Replacement
}

type Main struct {
	PushbulletToken  string
	PushbulletDevice string
	BaseFolder       string
	Series           []Series
}

// Parse will read, parse and validate config supplied in reader
func Parse(r io.Reader) (config *Main, err error) {
	config = &Main{}
	_, err = toml.DecodeReader(r, config)

	if err != nil {
		return nil, err
	}

	// Validate the config
	if config.BaseFolder == "" {
		return nil, errors.New("No BaseFolder set")
	}

	for _, serie := range config.Series {
		if serie.Key == "" {
			return nil, errors.New("No series key set for series")
		}

		if serie.Series != nil {
			_, err := regexp.Compile(serie.Series.Regex)
			if err != nil {
				return nil, errors.New("Failed to compile regex")
			}
		}

		if serie.Season != nil {
			_, err := regexp.Compile(serie.Season.Regex)
			if err != nil {
				return nil, errors.New("Failed to compile regex")
			}
		}

		if serie.Episode != nil {
			_, err := regexp.Compile(serie.Episode.Regex)
			if err != nil {
				return nil, errors.New("Failed to compile regex")
			}
		}
	}

	return
}
