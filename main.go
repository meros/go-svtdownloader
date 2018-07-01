package main

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/meros/go-svtdownloader/libs/confreader"
	"github.com/meros/go-svtdownloader/libs/epdownloader"
	"github.com/meros/go-svtdownloader/libs/eplister"
	"github.com/meros/go-svtdownloader/libs/epnamer"
)

func main() {
	config := kingpin.Flag("config", "Config file").Short('c').Required().String()
	forever := kingpin.Flag("forever", "Keep running forever").Short('f').Bool()

	kingpin.Parse()

	file, err := os.Open(*config)
	if err != nil {
		log.Fatal("Failed to open config file: ", err)
	}

	mainConfig, err := confreader.Parse(file)
	if err != nil {
		log.Fatal("Failed to parse config file: ", err)
	}

	mainConfig.Notifier.Notify("Starting up server", "")

	for {
		for _, serie := range mainConfig.Series {
			log.Println("Fetching series", serie.Key)
			eps, err := eplister.Get(serie.Key)
			if err != nil {
				log.Fatal("Failed to fetch serie", err)
			}

			epnamerOptions := epnamer.Options{
				Series: &epnamer.Replacement{
					Re:          serie.SeriesRegexp,
					Replacement: serie.SeriesReplacement},
				Season: &epnamer.Replacement{
					Re:          serie.SeasonRegexp,
					Replacement: serie.SeasonReplacement},
				Episode: &epnamer.Replacement{
					Re:          serie.EpisodeRegexp,
					Replacement: serie.EpisodeReplacement},
				Template: serie.FilenameTemplate}

			for _, ep := range eps {
				filename, _ := epnamer.Filename(ep, epnamerOptions)
				filename = path.Join(mainConfig.BaseFolder, filename)

				err := epdownloader.Get(ep, filename)
				if err == nil {
					mainConfig.Notifier.Notify("Episode downloaded", ep.Series+" "+ep.Season+" "+ep.Episode+" has been downloaded")
				}
			}
		}

		if !(*forever) {
			return
		}

		log.Println("Sleeping for 10 minutes and checking again at that point")
		time.Sleep(10 * time.Minute)

	}
}
