package confreader

import (
	"errors"
	"io"
	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/alecthomas/template"
	pushbullet "github.com/xconstruct/go-pushbullet"
)

type tomlSerie struct {
	Key                string
	FilenameTemplate   string
	SeriesRegexp       string
	SeriesReplacement  string
	SeasonRegexp       string
	SeasonReplacement  string
	EpisodeRegexp      string
	EpisodeReplacement string
}

type tomlPushBullet struct {
	Token  string
	Device string
}

type tomlMain struct {
	PushBullet *tomlPushBullet
	BaseFolder string
	Series     []tomlSerie
}

// Notifier is an interface used to send notifications to the user
type Notifier interface {
	Notify(title string, body string)
}

// Transformer is an interface used to transform a string into another string
type Transformer interface {
	Transform(text string) string
}

// ConfigSeries represents a series to download
type ConfigSeries struct {
	// Key is the path of the series as in https://www.svtplay.se/thunderbirds where thunderbirds is the key
	Key string
	// FilenameTemplate is a template that generates a file name
	// Should be fed a struct with keys Series, Season, Episode
	FilenameTemplate   *template.Template
	SeriesTransformer  Transformer
	SeasonTransformer  Transformer
	EpisodeTransformer Transformer
}

// ConfigMain is the full collection of configuration for the svtdownloader
type ConfigMain struct {
	BaseFolder string
	Notifier   Notifier
	Series     []ConfigSeries
}

type pushBulletNotifier struct {
	dev *pushbullet.Device
}

func (pbn *pushBulletNotifier) Notify(title string, body string) {
	pbn.dev.PushNote(title, body)
}

type emptyNotifier struct{}

func (en *emptyNotifier) Notify(title string, body string) {}

type regexpTransformer struct {
	re *regexp.Regexp
	rs string
}

func (ret *regexpTransformer) Transform(text string) string {
	return ret.re.ReplaceAllString(text, ret.rs)
}

// Parse will read, parse and validate config supplied in reader
func Parse(r io.Reader) (configMain *ConfigMain, err error) {
	configMain = &ConfigMain{}

	tomlMain := &tomlMain{}
	metadata, err := toml.DecodeReader(r, tomlMain)

	// Validate toml
	if err != nil {
		return nil, err
	}

	if len(metadata.Undecoded()) > 0 {
		return nil, errors.New("Unexpected keys found in config")
	}

	// Validate the config
	if tomlMain.BaseFolder == "" {
		return nil, errors.New("No basefolder set, set to output path for all series")
	}
	configMain.BaseFolder = tomlMain.BaseFolder

	if tomlMain.PushBullet != nil {
		if tomlMain.PushBullet.Token == "" || tomlMain.PushBullet.Device == "" {
			return nil, errors.New("Either token or device was empty string for pushbullet config")
		}

		pb := pushbullet.New(tomlMain.PushBullet.Token)
		dev, err := pb.Device(tomlMain.PushBullet.Device)

		if err != nil {
			return nil, err
		}

		configMain.Notifier = &pushBulletNotifier{dev}
	} else {
		configMain.Notifier = &emptyNotifier{}
	}

	for _, tomlSerie := range tomlMain.Series {
		configSeries := ConfigSeries{}

		if tomlSerie.Key == "" {
			return nil, errors.New("No series key set for series")
		}
		configSeries.Key = tomlSerie.Key

		if tomlSerie.FilenameTemplate == "" {
			return nil, errors.New("No series filenametemplate set for series")
		}
		configSeries.FilenameTemplate, err = template.New("FilenameTemplate").Parse(tomlSerie.FilenameTemplate)
		if err != nil {
			return nil, err
		}

		if tomlSerie.SeriesRegexp != "" {
			re, err := regexp.Compile(tomlSerie.SeriesRegexp)
			if err != nil {
				return nil, err
			}

			configSeries.SeriesTransformer = &regexpTransformer{
				re,
				tomlSerie.SeriesReplacement}
		} else {
			configSeries.SeriesTransformer = &regexpTransformer{
				regexp.MustCompile(".*"),
				tomlSerie.SeriesReplacement}
		}

		if tomlSerie.SeasonRegexp != "" {
			re, err := regexp.Compile(tomlSerie.SeasonRegexp)
			if err != nil {
				return nil, err
			}

			configSeries.SeasonTransformer = &regexpTransformer{
				re,
				tomlSerie.SeasonReplacement}
		} else {
			configSeries.SeasonTransformer = &regexpTransformer{
				regexp.MustCompile(".*"),
				tomlSerie.SeasonReplacement}
		}

		if tomlSerie.EpisodeRegexp != "" {
			re, err := regexp.Compile(tomlSerie.EpisodeRegexp)
			if err != nil {
				return nil, err
			}

			configSeries.EpisodeTransformer = &regexpTransformer{
				re,
				tomlSerie.EpisodeReplacement}
		} else {
			configSeries.EpisodeTransformer = &regexpTransformer{
				regexp.MustCompile(".*"),
				tomlSerie.EpisodeReplacement}
		}

		configMain.Series = append(configMain.Series, configSeries)
	}

	return configMain, nil
}
